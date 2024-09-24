package request

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/request"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/curescan"
)

type CreateTemplate struct {
	TemplateType    string `json:"templateType"`
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
