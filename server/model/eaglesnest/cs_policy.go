package eaglesnest

import (
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/lib/pq"
)

type Policy struct {
	global.GvaModel
	global.CsModel
	PolicyName     string         `json:"policyName" gorm:"column:policy_name;type:text;index;not null;comment:策略名称"`
	PolicyDesc     string         `json:"policyDesc" gorm:"column:policy_desc;type:text;comment:策略描述"`
	ScanType       pq.StringArray `json:"scanType" gorm:"column:scan_type;type:text[];comment:扫描类型"`
	PolicyConfig   string         `json:"policyConfig" gorm:"column:policy_config;type:text;comment:策略配置"`
	OnlineCheck    bool           `json:"onlineCheck" gorm:"column:online_check;type:bool;comment:在线检测"`
	OnlineConfig   string         `json:"onlineConfig" gorm:"column:online_config;type:text;comment:在线检测配置"`
	PortScan       bool           `json:"portScan" gorm:"column:port_scan;type:bool;comment:端口扫描"`
	PortScanConfig string         `json:"portScanConfig" gorm:"column:port_scan_config;type:text;comment:端口扫描配置"`
	Templates      pq.Int64Array  `json:"templates" gorm:"column:templates;type:int8[];comment:模板"`
	IgnoredIP      pq.StringArray `json:"ignoredIp" gorm:"column:ignored_ip;type:text[];comment:忽略的IP"`
}

func (Policy) TableName() string {
	return "cs_policy"
}
