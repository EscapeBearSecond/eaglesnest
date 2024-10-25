package eaglesnest

import (
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/lib/pq"
)

type Asset struct {
	global.GvaModel
	CreatedBy    uint          `gorm:"column:created_by;type:int8;uniqueIndex:idx_ip_created_by;comment:创建者" json:"CreatedBy"`
	UpdatedBy    uint          `gorm:"column:updated_by;type:int8;comment:更新者" json:"UpdatedBy"`
	AssetName    string        `json:"assetName" gorm:"type:text;column:asset_name;comment:资产名称"`
	AssetIP      string        `json:"assetIp" gorm:"type:text;uniqueIndex:idx_ip_created_by;column:asset_ip;comment:资产IP"`
	AreaName     string        `json:"areaName" gorm:"type:text;column:area_name;comment:资产所属区域名称"`
	AssetArea    uint          `json:"assetArea" gorm:"type:int8;column:asset_area;comment:资产所属区域"`
	AssetType    string        `json:"assetType" gorm:"type:text;column:asset_type;comment:资产类型"`
	OpenPorts    pq.Int64Array `json:"openPorts" gorm:"type:int8[];column:open_ports;comment:开放端口"`
	SystemType   string        `json:"systemType" gorm:"type:text;column:system_type;comment:系统类型"`
	Manufacturer string        `json:"manufacturer" gorm:"type:text;column:manufacturer;comment:厂商"`
	AssetModel   string        `json:"assetModel" gorm:"type:text;column:asset_model;comment:型号"`
	TTL          int64         `json:"ttl" gorm:"type:int;column:ttl;comment:TTL值"`
}

func (Asset) TableName() string {
	return "cs_asset"
}
