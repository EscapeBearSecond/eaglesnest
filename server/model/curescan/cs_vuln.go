package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONB map[string]interface{}
type JSONArray []string
type Vuln struct {
	global.GvaModel
	TemplateID     string    `gorm:"column:template_id;type:text;not null;comment:模板id" json:"templateId"`
	Name           string    `gorm:"column:name;type:text;not null;comment:漏洞名称" json:"name"`
	Author         string    `gorm:"column:author;type:text;not null;comment:漏洞作者" json:"author"`
	Severity       string    `gorm:"column:severity;type:text;not null;comment:漏洞等级" json:"severity"`
	Description    string    `gorm:"column:description;type:text;comment:漏洞描述" json:"description"`
	Reference      JSONArray `gorm:"column:reference;type:json;comment:引用信息" json:"reference"`
	Classification JSONB     `gorm:"column:classification;type:jsonb;comment:其他分类信息" json:"classification"`
	Remediation    string    `gorm:"column:remediation;type:text;comment:修复建议" json:"remediation"`
}

// 实现 Scanner 接口，将数据库中的 jsonb 转换为 Go 结构
func (sa *JSONArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, sa)
}

// 实现 Valuer 接口，将 Go 结构转换为 jsonb 存入数据库
func (sa JSONArray) Value() (driver.Value, error) {
	return json.Marshal(sa)
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
