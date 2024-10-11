package request

import (
	"github.com/EscapeBearSecond/curescan/server/model/common/request"
	"github.com/EscapeBearSecond/curescan/server/model/curescan"
)

type SearchJobResult struct {
	curescan.JobResultItem
	request.PageInfo
	DistinctFields []string `json:"distinctFields"`
	OrderKey       string   `json:"orderKey"` // 排序
	Desc           bool     `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
