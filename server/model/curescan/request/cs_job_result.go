package request

import (
	"47.103.136.241/goprojects/curescan/server/model/common/request"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
)

type SearchJobResult struct {
	curescan.JobResultItem
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
