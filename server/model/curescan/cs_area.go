package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"github.com/lib/pq"
)

// Area 区域
type Area struct {
	global.GvaModel
	global.CsModel
	AreaName string         `json:"areaName" gorm:"type:varchar(20);not null;index;column:area_name;"`
	AreaIP   pq.StringArray `json:"areaIp" gorm:"type:text[];not null;column:area_ip;"`
	AreaDesc string         `json:"areaDesc" gorm:"type:varchar(100);column:area_desc;"`
}

func (Area) TableName() string {
	return "cs_area"
}
