package eaglesnest

import (
	"errors"
	"github.com/EscapeBearSecond/eaglesnest/server/core/meta"
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/falcon/pkg/license"
	"gorm.io/gorm"
	"time"
)

type SystemInfoService struct {
}

func (s *SystemInfoService) InitSystemInfo() error {
	// 读取../version.ini文件
	systemVersion := "0.0.0"

	vulnVersion := "0.0.0"
	lastUpdateDate := time.Now().Format("2006-01-02 15:04:05")
	watcher, err := license.Watch("./license.json")
	if err != nil {
		return err
	}
	defer watcher.Stop()

	licenseExpiration := license.L().ExpiresAt
	var systemInfo = eaglesnest.SystemInfo{}
	err = global.GVA_DB.First(&systemInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有记录，创建一条新的记录
		systemInfo = eaglesnest.SystemInfo{
			LastUpdateDate:    lastUpdateDate,
			SystemVersion:     systemVersion,
			VulnVersion:       vulnVersion,
			LicenseExpiration: licenseExpiration,
		}
		// 使用 Create 方法插入新记录
		return global.GVA_DB.Create(&systemInfo).Error
	} else if err != nil {
		return err
	} else {
		// 如果找到了记录，更新现有的记录
		systemInfo.LastUpdateDate = lastUpdateDate
		systemInfo.SystemVersion = systemVersion
		systemInfo.VulnVersion = vulnVersion
		systemInfo.LicenseExpiration = licenseExpiration
	}
	return global.GVA_DB.Save(&systemInfo).Error
}

func (s *SystemInfoService) GetSystemInfo() (*eaglesnest.SystemInfo, error) {
	var systemInfo eaglesnest.SystemInfo
	err := global.GVA_DB.First(&systemInfo).Error
	if err != nil {
		return nil, err
	}
	systemInfo.SystemVersion = meta.BuildVer
	return &systemInfo, err
}

func (s *SystemInfoService) UpdateSystemInfo(info *eaglesnest.SystemInfo) error {
	return global.GVA_DB.Save(info).Error
}
