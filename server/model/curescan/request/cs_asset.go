package request

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/request"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
)

type SearchAsset struct {
	*curescan.Asset
	request.PageInfo
	global.CsModel
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
