package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
)

type Template struct {
	global.GvaModel
	TemplateName    string `gorm:"column:template_name;type:varchar(20);uniqueIndex;not null;comment:模板名称" json:"templateName"`
	TemplateType    uint   `gorm:"column:template_type;type:int8;not null;comment:模板类型" json:"templateType"`
	TemplateDesc    string `gorm:"column:template_desc;type:varchar(100);not null;comment:模板描述" json:"templateDesc"`
	TemplateContent string `gorm:"column:template_content;type:text;not null;comment:模板内容" json:"templateContent"`
}

func (Template) TableName() string {
	return "cs_template"
}
