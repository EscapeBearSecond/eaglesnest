package curescan

import "47.103.136.241/goprojects/curescan/server/global"

type JobResultItem struct {
	global.GvaModel
	global.CsModel
	Name             string   `json:"name" gorm:"type:text;column:name;comment:任务名称"`
	Kind             string   `json:"kind" gorm:"type:text;column:kind;comment:任务类型"`
	TemplateID       string   `json:"templateId" gorm:"type:text;column:template_id;comment:模板ID"`
	TemplateName     string   `json:"templateName" gorm:"type:text;column:template_name;comment:模板名称"`
	Type             string   `json:"type" gorm:"type:text;column:type;comment:模板类型"`
	Severity         string   `json:"severity" gorm:"type:text;column:severity;comment:模板等级"`
	Host             string   `json:"host" gorm:"type:text;column:host;comment:主机"`
	Port             string   `json:"port" gorm:"type:text;column:port;comment:端口"`
	Scheme           string   `json:"scheme" gorm:"type:text;column:scheme;comment:协议"`
	URL              string   `json:"url" gorm:"type:text;column:url;comment:url,unique"`
	Path             string   `json:"path" gorm:"type:text;column:path;comment:路径"`
	Matched          string   `json:"matched" gorm:"type:text;column:matched;comment:匹配结果"`
	ExtractedResults []string `json:"extracted_results" gorm:"type:text[];column:extracted_results;comment:提取结果"`
	Description      string   `json:"description" gorm:"type:text;column:description;comment:描述"`
	TaskID           uint     `json:"taskId" gorm:"type:int8;column:task_id;comment:任务ID"`
	EntryID          string   `json:"entryId" gorm:"type:text;column:entry_id;comment:条目ID"`
	Remediation      string   `json:"remediation" gorm:"type:text;column:remediation;comment:修复建议"`
	Tag1             string   `json:"tag1" gorm:"type:text;column:tag1;"`
	Tag2             string   `json:"tag2" gorm:"type:text;column:tag2;"`
	Tag3             string   `json:"tag3" gorm:"type:text;column:tag3;"`
	Tag4             string   `json:"tag4" gorm:"type:text;column:tag4;"`
}

func (JobResultItem) TableName() string {
	return "cs_job_result"
}
