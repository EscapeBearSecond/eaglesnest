package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/request"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AreaService struct {
}

// CreateArea 创建一个新的区域，不允许有重复的区域名称。
func (a *AreaService) CreateArea(area *curescan.Area) error {
	if !errors.Is(global.GVA_DB.Select("area_name").First(&curescan.Area{}, "area_name=?", area.AreaName).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同区域名称，不允许创建")
	}

	return global.GVA_DB.Create(&area).Error
}

// DeleteArea 根据区域ID删除区域，该删除是逻辑删除，通过将deleted_at字段置为删除时间。
func (a *AreaService) DeleteArea(id int) error {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&curescan.Area{}, id).Error
		if err != nil {
			return err
		}
		err = tx.Model(&curescan.Asset{}).Where("asset_area = ?", id).Update("area_name", "未知").Error
		return err
	})
	return err
}

// UpdateArea 更新区域信息，更新后的区域名称不允许重复。
func (a *AreaService) UpdateArea(area *curescan.Area) error {
	var existingRecord curescan.Area
	err := global.GVA_DB.Select("id", "area_name").Where("area_name=?", area.AreaName).First(&existingRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.GVA_DB.Save(&area).Error
		}
		return err
	}
	if existingRecord.ID != area.ID {
		return errors.New("区域名称已被占用，不允许修改")
	}
	return global.GVA_DB.Save(&area).Error
}

// GetAreaById 根据区域ID获取区域详情。
func (a *AreaService) GetAreaById(id int) (*curescan.Area, error) {
	var area curescan.Area
	err := global.GVA_DB.Select("id", "area_name", "area_desc", "area_ip",
		"created_at", "updated_at", "deleted_at").Where("id=?", id).First(&area).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	return &area, nil
}

// GetAreaList 获取区域列表，该方法会根据页码信息和排序信息返回分页后的区域信息。调用该方法需要传递的参数有4个，第一个为过滤信息，也就是要查询的区域信息或关键字；
// 第二个参数是分页信息；第三个参数是排序字段，第四个参数是是否倒序。如查询区域名称为“南京”，且要按照ID字段倒序排序，则参数 area.AreaName="南京", page.Page=1,
// page.PageInfo=10, order="id", desc=true
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
	if order != "" {
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
