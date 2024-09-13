package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	request2 "47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type TemplateService struct {
}

// CreateTemplate 创建模板，不允许重复的模板名称。
func (t *TemplateService) CreateTemplate(template *curescan.Template) error {
	if !errors.Is(global.GVA_DB.Select("template_name").First(&curescan.Template{}, "template_name = ?", template.TemplateName).Error, gorm.ErrRecordNotFound) {
		return errors.New("模板已存在，请查看模板名是否正确")
	}
	return global.GVA_DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "template_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"template_desc",
			"template_content",
			"tag1",
			"tag2",
			"tag3",
			"tag4",
			"template_type",
			"template_name",
			"deleted_at",
			"updated_at", // 确保更新操作时更新时间戳
			"updated_by", // 确保更新操作时更新操作用户
		}),
	}).Create(template).Error

}

// DeleteTemplate 根据模板ID删除模板
func (t *TemplateService) DeleteTemplate(id int) error {
	return global.GVA_DB.Delete(&curescan.Template{}, id).Error
}

// GetTemplateById 根据模板ID获取模板详情信息，包括模板内容
func (t *TemplateService) GetTemplateById(id int) (*curescan.Template, error) {
	var template curescan.Template
	err := global.GVA_DB.Select("id", "template_name", "template_type", "template_id", "template_desc", "template_content",
		"created_at", "updated_at", "deleted_at", "tag1", "tag2", "tag3", "tag4").Where("id = ?", id).First(&template).Error
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
	err := global.GVA_DB.Select("id", "template_name", "template_type", "template_id", "template_desc", "template_content",
		"created_at", "updated_at", "deleted_at").Where("id in (?)", ids).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// GetTemplateList 获取模板列表，该方法返回除模板内容外的所有信息。如果想要获取模板内容，需要调用GetTemplateById方法。
func (t *TemplateService) GetTemplateList(searchTemplate request2.SearchTemplate) (list []*curescan.Template, total int64, err error) {
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
	var templates []*curescan.Template

	if template.TemplateType != "" {
		db = db.Where("template_type = ?", template.TemplateType)
	}
	if template.TemplateId != "" {
		db = db.Where("template_id = ?", template.TemplateId)
	}
	if template.Tag1 != "" && template.Tag1 != "''" {
		db = db.Where("tag1 = ?", template.Tag1)
	}
	if template.Tag2 != "" && template.Tag2 != "''" {
		db = db.Where("tag2 = ?", template.Tag2)
	}
	if template.Tag3 != "" && template.Tag3 != "''" {
		db = db.Where("tag3 = ?", template.Tag3)
	}
	if template.Tag4 != "" && template.Tag4 != "''" {
		db = db.Where("tag4 = ?", template.Tag4)
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
		orderMap["template_id"] = true
		orderMap["tag1"] = true
		orderMap["tag2"] = true
		orderMap["tag3"] = true
		orderMap["tag4"] = true
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
	uniqueTemplates := make(map[string]*curescan.Template)
	for _, template := range templates {
		if _, exists := uniqueTemplates[template.TemplateId]; !exists {
			uniqueTemplates[template.TemplateId] = template
		}
	}

	var deduplicatedTemplates []*curescan.Template
	for _, template := range uniqueTemplates {
		deduplicatedTemplates = append(deduplicatedTemplates, template)
	}

	// 开启事务，确保批量操作的原子性
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 执行批量插入，处理冲突
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "template_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"template_desc",
				"template_content",
				"tag1",
				"tag2",
				"tag3",
				"tag4",
				"template_type",
				"template_name",
				"deleted_at",
				"updated_at", // 确保更新操作时更新时间戳
				"updated_by", // 确保更新操作时更新操作用户
			}),
		}).CreateInBatches(deduplicatedTemplates, 100).Error; err != nil {
			return err // 如果发生错误，回滚事务
		}
		return nil // 提交事务
	})
}

func (t *TemplateService) ParseTemplateContent(template *curescan.Template) (err error) {
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer([]byte(template.TemplateContent)))
	if err != nil {
		err = errors.New("模板内容有误, 请检查模板内容是否是yaml格式")
		return
	}
	tagsStr := viper.GetString("info.tags")
	templateName := viper.GetString("info.name")

	templateId := viper.GetString("id")
	if tagsStr == "" || templateName == "" || templateId == "" {
		err = errors.New("模板内容有误, 请检查模板名称, 模板id, 模板标签是否填写完整")
		return
	}
	template.TemplateName = templateName
	template.TemplateDesc = viper.GetString("info.description")
	template.TemplateId = templateId
	tags := strings.Split(tagsStr, ",")

	if len(tags) == 1 {
		template.Tag1 = tags[0]
		template.Tag2 = "未知"
		template.Tag3 = "未知"
		template.Tag4 = "未知"
	}
	if len(tags) == 2 {
		template.Tag1 = tags[0]
		template.Tag2 = tags[1]
		template.Tag3 = "未知"
		template.Tag4 = "未知"
	}
	if len(tags) == 3 {
		template.Tag2 = tags[1]
		template.Tag3 = tags[2]
		template.Tag1 = tags[0]
		template.Tag4 = "未知"
	}
	if len(tags) == 4 {
		template.Tag1 = tags[0]
		template.Tag2 = tags[1]
		template.Tag3 = tags[2]
		template.Tag4 = tags[3]
	}
	return
}
