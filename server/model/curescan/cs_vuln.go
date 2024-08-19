package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"github.com/lib/pq"
)

type Vuln struct {
	global.GvaModel
	TemplateID     string                 `gorm:"column:template_id;type:text;uniqueIndex;not null;comment:模板id" json:"templateId"`
	Name           string                 `gorm:"column:name;type:text;uniqueIndex;not null;comment:漏洞名称" json:"name"`
	Author         string                 `gorm:"column:author;type:text;not null;comment:漏洞作者" json:"author"`
	Severity       string                 `gorm:"column:severity;type:text;not null;comment:漏洞等级" json:"severity"`
	Description    string                 `gorm:"column:description;type:text;not null;comment:漏洞描述" json:"description"`
	Reference      pq.StringArray         `gorm:"column:reference;type:text[];not null;comment:引用信息" json:"reference"`
	Classification map[string]interface{} `gorm:"column:classification;type:jsonb;not null;comment:其他分类信息" json:"classification"`
	Remediation    string                 `gorm:"column:remediation;type:text;not null;comment:修复建议" json:"remediation"`
}

func (v *Vuln) TableName() string {
	return "cs_vuln"
}
