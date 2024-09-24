package request

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/request"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/system"
)

// api分页条件查询及排序结构体
type SearchApiParams struct {
	system.SysApi
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
