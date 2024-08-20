package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

type JSONB map[string]interface{}
type Vuln struct {
	global.GvaModel
	TemplateID     string         `gorm:"column:template_id;type:text;uniqueIndex;not null;comment:模板id" json:"templateId"`
	Name           string         `gorm:"column:name;type:text;uniqueIndex;not null;comment:漏洞名称" json:"name"`
	Author         string         `gorm:"column:author;type:text;not null;comment:漏洞作者" json:"author"`
	Severity       string         `gorm:"column:severity;type:text;not null;comment:漏洞等级" json:"severity"`
	Description    string         `gorm:"column:description;type:text;comment:漏洞描述" json:"description"`
	Reference      pq.StringArray `gorm:"column:reference;type:text[];comment:引用信息" json:"reference"`
	Classification JSONB          `gorm:"column:classification;type:jsonb;comment:其他分类信息" json:"classification"`
	Remediation    string         `gorm:"column:remediation;type:text;comment:修复建议" json:"remediation"`
}

func (Vuln) TableName() string {
	return "cs_vuln"
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

func (j *JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}
