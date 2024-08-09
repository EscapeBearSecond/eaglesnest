package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"fmt"
	"gorm.io/gorm"
)

type VulnService struct {
}

func (s *VulnService) GetVulnList(searchVuln *request.SearchVuln) (list interface{}, total int64, err error) {
	vuln := searchVuln.Vuln
	page := searchVuln.PageInfo
	order := searchVuln.OrderKey
	desc := searchVuln.Desc
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)

	var db *gorm.DB
	db = global.GVA_DB.Model(&curescan.Template{})
	var templates []curescan.Template
	if vuln.Name != "" {
		db = db.Where("name LIKE", "%"+vuln.Name+"%")
	}
	if vuln.Author != "" {
		db = db.Where("author LIKE ?", "%"+vuln.Author+"%")
	}
	if vuln.Severity != "" {
		db = db.Where("severity LIKE ?", "%"+vuln.Severity+"%")
	}
	if vuln.Reference != "" {
		db = db.Where("reference LIKE ?", "%"+vuln.Reference+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return templates, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool, 3)
		orderMap["id"] = true
		orderMap["name"] = true
		orderMap["author"] = true
		orderMap["severity"] = true
		orderMap["reference"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return templates, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&templates).Error
	return templates, total, err

}
