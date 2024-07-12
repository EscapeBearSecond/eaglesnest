package curescan

import (
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"errors"
	"fmt"

	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"gorm.io/gorm"
)

type PolicyService struct {
}

// CreatePolicy 创建策略, 不允许有重复的策略名称.
func (p *PolicyService) CreatePolicy(policy *curescan.Policy) error {
	if !errors.Is(global.GVA_DB.Select("policy_name").First(&curescan.Policy{}, "policy_name = ?", policy.PolicyName).Error, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("存在相同策略名称【%s】,不允许创建", policy.PolicyName))
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
		return errors.New(fmt.Sprintf("名称为【%s】的策略已存在，不允许修改", policy.PolicyName))
	}
	return global.GVA_DB.Save(&policy).Error
}

// DeletePolicy 根据策略ID删除策略, 该删除是逻辑删除, 通过将deleted_at字段的置为当前时间来进行删除.
func (p *PolicyService) DeletePolicy(id int) error {
	return global.GVA_DB.Delete(&curescan.Policy{}, id).Error
}

// GetPolicyById 根据策略ID获取策略详情信息, 该方法会返回一条策略的全部信息.
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

// GetPolicyList 获取策略列表. 该方法会根据页码信息和排序信息返回分页后的策略信息.
// 对于每一条策略, 该方法不会返回全部的信息, 只会返回内容较少且常用的信息, 如果想要获取一条策略的详细信息, 需要调用GetPolicyById方法
// 调用该方法需要传递的参数有4个, 第一个为过滤信息, 也就是要查询的策略信息或关键字;
// 第二个参数是分页信息; 第三个参数是排序字段; 第四个参数是是否倒序. 如查询扫描类型为"漏洞扫描", 且要按照策略名称字段倒序排序, 则参数 policy.ScanType=["漏洞扫描"], page.Page=1,
// page.PageInfo=10, order="policyName", desc=true
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
