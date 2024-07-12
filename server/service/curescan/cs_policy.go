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

func (p *PolicyService) CreatePolicy(policy *curescan.Policy) error {
	if !errors.Is(global.GVA_DB.Select("policy_name").First(&curescan.Policy{}, "policy_name = ?", policy.PolicyName).Error, gorm.ErrRecordNotFound) {
		return errors.New(fmt.Sprintf("存在相同策略名称【%s】,不允许创建", policy.PolicyName))
	}
	return global.GVA_DB.Create(policy).Error
}

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

func (p *PolicyService) DeletePolicy(id int) error {
	return global.GVA_DB.Delete(&curescan.Policy{}, id).Error
}

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
