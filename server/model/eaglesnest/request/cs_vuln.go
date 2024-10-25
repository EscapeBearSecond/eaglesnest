package request

import (
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
)

type SearchVuln struct {
	eaglesnest.Vuln
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
