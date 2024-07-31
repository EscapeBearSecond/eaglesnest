package curescan

import (
	"errors"
	"fmt"

	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	request2 "47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"gorm.io/gorm"
)

type TemplateService struct {
}

// CreateTemplate 创建模板，不允许重复的模板名称。
func (t *TemplateService) CreateTemplate(template *curescan.Template) error {
	if !errors.Is(global.GVA_DB.Select("template_name").First(&curescan.Template{}, "template_name = ?", template.TemplateName).Error, gorm.ErrRecordNotFound) {
		return errors.New("模板已存在，请查看模板名是否正确")
	}

	return global.GVA_DB.Create(template).Error
}

// DeleteTemplate 根据模板ID删除模板
func (t *TemplateService) DeleteTemplate(id int) error {
	return global.GVA_DB.Delete(&curescan.Template{}, id).Error
}

// GetTemplateById 根据模板ID获取模板详情信息，包括模板内容
func (t *TemplateService) GetTemplateById(id int) (*curescan.Template, error) {
	var template curescan.Template
	err := global.GVA_DB.Select("id", "template_name", "template_type", "template_desc", "template_content",
		"created_at", "updated_at", "deleted_at").Where("id = ?", id).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	return &template, nil
}

func (t *TemplateService) GetTemplatesByIds(ids []int64) ([]*curescan.Template, error) {
	var templates []*curescan.Template
	err := global.GVA_DB.Select("id", "template_name", "template_type", "template_desc", "template_content",
		"created_at", "updated_at", "deleted_at").Where("id in (?)", ids).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// GetTemplateList 获取模板列表，该方法返回除模板内容外的所有信息。如果想要获取模板内容，需要调用GetTemplateById方法。
func (t *TemplateService) GetTemplateList(searchTemplate request2.SearchTemplate) (list interface{}, total int64, err error) {
	template := searchTemplate.Template
	page := searchTemplate.PageInfo
	order := searchTemplate.OrderKey
	desc := searchTemplate.Desc
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	var db *gorm.DB
	if searchTemplate.IsAll {
		db = global.GVA_DB.Model(&curescan.Template{})
	} else {
		db = global.GVA_DB.Model(&curescan.Template{}).Omit("template_content")
	}
	var templates []curescan.Template
	if template.TemplateType != 0 {
		db = db.Where("template_type = ?", template.TemplateType)
	}
	if template.TemplateName != "" {
		db = db.Where("template_name LIKE ?", "%"+template.TemplateName+"%")
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
		orderMap["template_name"] = true
		orderMap["template_type"] = true
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

// UpdateTemplate 更新模板信息，更新后的模板名称不允许重复
func (t *TemplateService) UpdateTemplate(template *curescan.Template) error {
	var existingRecord curescan.Template
	err := global.GVA_DB.Select("id", "template_name").Where("template_name=?", template.TemplateName).First(&existingRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.GVA_DB.Save(&template).Error
		}
		return err
	}
	if existingRecord.ID != template.ID {
		return errors.New("模板名称已被占用，不允许修改")
	}
	return global.GVA_DB.Save(&template).Error
}

func (t *TemplateService) BatchAdd(templates []*curescan.Template) error {
	return global.GVA_DB.Model(&curescan.Template{}).CreateInBatches(templates, 100).Error
}
