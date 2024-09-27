package curescan

type SystemInfo struct {
	ID                uint   `gorm:"primaryKey" json:"ID"` // 主键ID
	SystemVersion     string `json:"systemVersion" gorm:"type:text;column:system_version;comment:系統版本;"`
	VulnVersion       string `json:"vulnVersion" gorm:"type:text;column:vuln_version;comment:漏洞库版本"`
	LicenseExpiration string `json:"licenseExpiration" gorm:"type:text;column:license_expiration;comment:证书过期时间"`
	LastUpdateDate    string `json:"lastUpdateDate" gorm:"type:text;column:last_update_date;comment:最新更新时间"`
}

func (SystemInfo) TableName() string {
	return "cs_system_info"
}
