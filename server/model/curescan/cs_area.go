package curescan

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	"github.com/lib/pq"
)

// Area 区域
type Area struct {
	global.GvaModel
	global.CsModel
	AreaName string         `json:"areaName" gorm:"type:text;not null;index;column:area_name;"`
	AreaIP   pq.StringArray `json:"areaIp" gorm:"type:text[];not null;column:area_ip;"`
	AreaDesc string         `json:"areaDesc" gorm:"type:text;column:area_desc;"`
}

func (Area) TableName() string {
	return "cs_area"
}
