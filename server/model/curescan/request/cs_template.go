package request

import (
	"github.com/EscapeBearSecond/curescan/server/model/common/request"
	"github.com/EscapeBearSecond/curescan/server/model/curescan"
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
