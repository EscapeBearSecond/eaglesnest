package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"github.com/lib/pq"
)

// Area 区域
type Area struct {
	global.GvaModel
	AreaName string         `json:"areaName" gorm:"type:varchar(20);not null;index;column:area_name;"`
	AreaIP   pq.StringArray `json:"areaIP" gorm:"type:text[];not null;column:area_ip;"`
	AreaDesc string         `json:"areaDesc" gorm:"type:text;column:area_desc;"`
}

func (Area) TableName() string {
	return "cs_area"
}
