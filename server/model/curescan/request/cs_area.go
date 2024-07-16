package request

import (
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
)

type CreateArea struct {
	AreaName string   `json:"areaName"`
	AreaIP   []string `json:"areaIp"`
	AreaDesc string   `json:"areaDesc"`
}

type UpdateArea struct {
	ID       uint     `json:"id"`
	AreaName string   `json:"areaName"`
	AreaIP   []string `json:"areaIp"`
	AreaDesc string   `json:"areaDesc"`
}

type SearchArea struct {
	curescan.Area
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
