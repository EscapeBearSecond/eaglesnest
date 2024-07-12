package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"github.com/lib/pq"
)

type Policy struct {
	global.GvaModel
	PolicyName     string         `gorm:"column:policy_name;type:varchar(50);uniqueIndex;not null;comment:策略名称"`
	PolicyDesc     string         `gorm:"column:policy_desc;type:varchar(100);comment:策略描述"`
	ScanType       pq.StringArray `gorm:"column:scan_type;type:text[];comment:扫描类型"`
	PolicyConfig   string         `gorm:"column:policy_config;type:text;comment:策略配置"`
	OnlineCheck    bool           `gorm:"column:online_check;type:bool;comment:在线检测"`
	OnlineConfig   string         `gorm:"column:online_config;type:text;comment:在线检测配置"`
	PortScan       bool           `gorm:"column:port_scan;type:bool;comment:端口扫描"`
	PortScanConfig string         `gorm:"column:port_scan_config;type:text;comment:端口扫描配置"`
	Templates      pq.Int64Array  `gorm:"column:templates;type:int8[];comment:模板"`
}

func (Policy) TableName() string {
	return "cs_policy"
}
