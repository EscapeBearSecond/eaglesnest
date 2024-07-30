package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"github.com/lib/pq"
)

type Asset struct {
	global.GvaModel
	AssetName    string        `json:"assetName" gorm:"type:varchar(255);column:asset_name;comment:资产名称"`
	AssetIP      string        `json:"assetIp" gorm:"type:varchar(255);column:asset_ip;comment:资产IP"`
	AssetArea    uint          `json:"assetArea" gorm:"type:int8;column:asset_area;comment:资产所属区域"`
	AssetType    string        `json:"assetType" gorm:"type:varchar(255);column:asset_type;comment:资产类型"`
	OpenPorts    pq.Int64Array `json:"openPorts" gorm:"type:int8[];column:open_ports;comment:开放端口"`
	SystemType   string        `json:"systemType" gorm:"type:varchar(255);column:system_type;comment:系统类型"`
	Manufacturer string        `json:"manufacturer" gorm:"type:varchar(255);column:manufacturer;comment:厂商"`
	AssetModel   string        `json:"assetModel" gorm:"type:varchar(255);column:asset_model;comment:型号"`
	TTL          int64         `json:"ttl" gorm:"type:int;column:ttl;comment:TTL值"`
}

func (Asset) TableName() string {
	return "cs_asset"
}
