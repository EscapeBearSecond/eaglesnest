package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type TemplateService struct {
}

func (t *TemplateService) CreateTemplate(template *curescan.Template) error {
	if !errors.Is(global.GVA_DB.First(&curescan.Template{}, "template_name = ?", template.TemplateName).Error, gorm.ErrRecordNotFound) {
		return errors.New("模板已存在，请查看模板名是否正确")
	}
	return global.GVA_DB.Create(template).Error
}

func (t *TemplateService) DeleteTemplate(id int) error {
	return global.GVA_DB.Delete(&curescan.Template{}, id).Error
}

func (t *TemplateService) GetTemplateById(id int) (*curescan.Template, error) {
	var template curescan.Template
	err := global.GVA_DB.Where("id = ?", id).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	return &template, nil
}

func (t *TemplateService) GetTemplateList(template curescan.Template, page request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.Template{}).Omit("template_content")
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

func (t *TemplateService) UpdateTemplate(template *curescan.Template) error {
	var existingRecord curescan.Template
	err := global.GVA_DB.Select("id", "template_name").Where("template_name=?", template.TemplateName).First(&existingRecord).Error
	if err != nil {
		return err
	}
	if existingRecord.ID != template.ID {
		return errors.New("模板名称已被占用，不允许修改")
	}
	return global.GVA_DB.Save(&template).Error
}
