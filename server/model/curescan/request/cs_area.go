package request

import (
	"github.com/EscapeBearSecond/curescan/server/model/common/request"
	"github.com/EscapeBearSecond/curescan/server/model/curescan"
)

type CreateArea struct {
	AreaName string   `json:"areaName"`
	AreaIP   []string `json:"areaIp" validate:"required,dive,ip_addr"`
	AreaDesc string   `json:"areaDesc"`
}

type UpdateArea struct {
	ID       uint     `json:"ID"`
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
