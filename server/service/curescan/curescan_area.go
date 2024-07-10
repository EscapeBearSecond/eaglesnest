package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"errors"
	"gorm.io/gorm"
)

type AreaService struct {
}

func (a *AreaService) CreateArea(area curescan.Area) error {
	if !errors.Is(global.GVA_DB.First(&curescan.Area{}, "area_name=?", area.AreaName).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同区域名称，不允许创建")
	}
	return global.GVA_DB.Create(&area).Error
}
