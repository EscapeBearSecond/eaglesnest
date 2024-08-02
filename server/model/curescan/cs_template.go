package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
)

type Template struct {
	global.GvaModel
	global.CsModel
	TemplateName    string `gorm:"column:template_name;type:text;uniqueIndex;not null;comment:模板名称" json:"templateName"`
	TemplateType    string `gorm:"column:template_type;type:text;not null;comment:模板类型" json:"templateType"`
	TemplateDesc    string `gorm:"column:template_desc;type:text;comment:模板描述" json:"templateDesc"`
	TemplateContent string `gorm:"column:template_content;type:text;not null;comment:模板内容" json:"templateContent"`
	TemplateId      string `gorm:"column:template_id;type:text;uniqueIndex;not null;comment:模板id" json:"templateId"`
	Tag1            string `gorm:"column:tag1;type:text;comment:tag1" json:"tag1"`
	Tag2            string `gorm:"column:tag2;type:text;comment:tag2" json:"tag2"`
	Tag3            string `gorm:"column:tag3;type:text;comment:tag3" json:"tag3"`
	Tag4            string `gorm:"column:tag4;type:text;comment:tag4" json:"tag4"`
}

func (Template) TableName() string {
	return "cs_template"
}
