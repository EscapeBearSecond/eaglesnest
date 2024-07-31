package request

import (
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
)

type CreateTemplate struct {
	TemplateName    string `json:"templateName"`
	TemplateType    uint   `json:"templateType,string"`
	TemplateDesc    string `json:"templateDesc"`
	TemplateContent string `json:"templateContent"`
}

type SearchTemplate struct {
	curescan.Template
	IsAll bool `json:"isAll"`
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

type UpdateTemplate struct {
	ID uint `json:"ID"`
	CreateTemplate
}
