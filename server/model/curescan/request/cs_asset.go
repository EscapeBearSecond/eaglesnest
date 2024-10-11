package request

import (
	"github.com/EscapeBearSecond/curescan/server/model/common/request"
	"github.com/EscapeBearSecond/curescan/server/model/curescan"
)

type SearchAsset struct {
	*curescan.Asset
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
