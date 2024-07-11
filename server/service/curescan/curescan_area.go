package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AreaService struct {
}

func (a *AreaService) CreateArea(area *curescan.Area) error {
	if !errors.Is(global.GVA_DB.First(&curescan.Area{}, "area_name=?", area.AreaName).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同区域名称，不允许创建")
	}

	return global.GVA_DB.Create(&area).Error
}

func (a *AreaService) DeleteArea(id int) error {
	return global.GVA_DB.Delete(&curescan.Area{}, id).Error
}

func (a *AreaService) UpdateArea(area *curescan.Area) error {
	var existingRecord curescan.Area
	err := global.GVA_DB.Where("area_name=?", area.AreaName).First(&existingRecord).Error
	if err != nil {
		return err
	}
	if existingRecord.ID != area.ID {
		return errors.New("区域名称已被占用，不允许修改")
	}
	return global.GVA_DB.Save(&area).Error
}

func (a *AreaService) GetAreaById(id int) (*curescan.Area, error) {
	var area curescan.Area
	err := global.GVA_DB.Where("id=?", id).First(&area).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	return &area, nil
}

func (a *AreaService) GetAreaList(area curescan.Area, page request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.Area{})
	var areas []curescan.Area
	if area.AreaName != "" {
		db = db.Where("area_name LIKE ?", "%"+area.AreaName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return areas, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if desc {
		orderMap := make(map[string]bool)
		orderMap["id"] = true
		orderMap["area_name"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return areas, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&areas).Error
	return areas, total, err
}
