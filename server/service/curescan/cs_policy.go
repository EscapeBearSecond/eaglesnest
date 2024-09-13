package curescan

import (
	"errors"
	"fmt"
	"math"

	request2 "47.103.136.241/goprojects/curescan/server/model/curescan/request"

	"47.103.136.241/goprojects/curescan/server/model/common/request"

	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"gorm.io/gorm"
)

type PolicyService struct {
}

var (
	taskService = &TaskService{}
)

// CreatePolicy 创建策略, 不允许有重复的策略名称.
func (p *PolicyService) CreatePolicy(policy *curescan.Policy) error {
	if !errors.Is(global.GVA_DB.Select("policy_name").First(&curescan.Policy{}, "policy_name = ?", policy.PolicyName).Error, gorm.ErrRecordNotFound) {
		return global.HasExisted
	}
	return global.GVA_DB.Create(policy).Error
}

// UpdatePolicy 更新策略信息, 更新后的策略名称不允许重复存在.
func (p *PolicyService) UpdatePolicy(policy *curescan.Policy) error {
	var existingRecord curescan.Policy
	err := global.GVA_DB.Select("id", "policy_name").Where("policy_name=?", policy.PolicyName).First(&existingRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.GVA_DB.Save(&policy).Error
		}
		return err
	}
	if existingRecord.ID != policy.ID {
		return fmt.Errorf("名称为【%s】的策略已存在，不允许修改", policy.PolicyName)
	}
	return global.GVA_DB.Save(&policy).Error
}

// DeletePolicy 根据策略ID删除策略, 该删除是逻辑删除, 通过将deleted_at字段的置为当前时间来进行删除.
func (p *PolicyService) DeletePolicy(id int) error {
	searchTask := request2.SearchTask{}
	searchTask.PageSize = math.MaxInt64
	searchTask.Page = 1
	searchTask.PolicyId = id
	_, total, err := taskService.GetTaskList(searchTask)
	if err != nil {
		return err
	}
	if total > 0 {
		return fmt.Errorf("策略存在关联任务，不允许删除，请先处理相关任务")
	}

	return global.GVA_DB.Delete(&curescan.Policy{}, id).Error
}

// GetPolicyById 根据策略ID获取策略信息。
// 该方法通过查询数据库，根据提供的ID获取特定的策略记录。
// 如果找到记录，它将返回该策略的详细信息；如果没有找到记录，它将返回一个错误。
//
// 参数:
//
//	id - 策略的唯一标识符。
//
// 返回值:
//
//	*curescan.Policy - 找到的策略对象，如果未找到则为nil。
//	error - 查询过程中发生的错误，如果未找到记录，则错误消息为“目标数据不存在”。
func (p *PolicyService) GetPolicyById(id int) (*curescan.Policy, error) {
	var policy curescan.Policy
	err := global.GVA_DB.Select("id", "policy_name", "policy_desc", "scan_type", "policy_config", "online_check",
		"online_config", "port_scan", "port_scan_config", "templates", "created_at", "updated_at", "deleted_at").Where("id = ?", id).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	return &policy, nil
}

// GetPolicyList 获取策略列表，根据给定的策略条件、分页信息和排序要求，返回策略列表、总条数和可能的错误。
// 该方法不会返回全部的信息，如果需要返回某条策略的全部信息，需要调用GetPolicyById方法。
//
// 参数:
//
//	policy: 策略对象，包含查询条件。
//	page: 分页信息，包含页码和每页大小。
//	order: 排序字段。
//	desc: 是否降序排序。
//
// 返回值:
//
//	list: 策略列表，类型为interface{}，具体为[]curescan.Policy。
//	total: 策略总条数。
//	err: 查询过程中可能出现的错误。
func (p *PolicyService) GetPolicyList(policy curescan.Policy, page request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.Policy{}).Omit("policy_config", "online_config", "port_scan_config")
	var policies []curescan.Policy
	if policy.PolicyName != "" {
		db = db.Where("policy_name LIKE ?", "%"+policy.PolicyName+"%")
	}
	if len(policy.ScanType) != 0 {
		db = db.Where("scan_type IN (?)", policy.ScanType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return policies, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool)
		orderMap["id"] = true
		orderMap["policy_name"] = true
		orderMap["scan_type"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return policies, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&policies).Error
	return policies, total, err
}
