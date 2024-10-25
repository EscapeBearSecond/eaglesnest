package eaglesnest

import (
	"errors"
	"fmt"

	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"gorm.io/gorm"
)

type AreaService struct {
}

// CreateArea 创建一个新的区域，不允许有重复的区域名称。
func (a *AreaService) CreateArea(area *eaglesnest.Area) error {
	if !errors.Is(global.GVA_DB.Select("area_name").First(&eaglesnest.Area{}, "area_name=?", area.AreaName).Error, gorm.ErrRecordNotFound) {
		return global.HasExisted
	}

	return global.GVA_DB.Create(&area).Error
}

// DeleteArea 根据区域ID删除区域，该删除是逻辑删除，通过将deleted_at字段置为删除时间。
func (a *AreaService) DeleteArea(id int) error {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&eaglesnest.Area{}, id).Error
		if err != nil {
			return err
		}
		err = tx.Model(&eaglesnest.Asset{}).Where("asset_area = ?", id).Update("area_name", "未知").Error
		return err
	})
	return err
}

// UpdateArea 更新区域信息，更新后的区域名称不允许重复。
func (a *AreaService) UpdateArea(area *eaglesnest.Area) error {
	var existingRecord eaglesnest.Area
	err := global.GVA_DB.Select("id", "area_name").Where("area_name=?", area.AreaName).First(&existingRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.GVA_DB.Save(&area).Error
		}
		return err
	}
	if existingRecord.ID != area.ID {
		return global.HasExisted
	}
	return global.GVA_DB.Save(&area).Error
}

// GetAreaById 根据区域ID获取区域详情。
func (a *AreaService) GetAreaById(id int) (*eaglesnest.Area, error) {
	var area eaglesnest.Area
	err := global.GVA_DB.Select("id", "area_name", "area_desc", "area_ip",
		"created_at", "updated_at", "deleted_at").Where("id=?", id).First(&area).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.NoDataFound
		}
		return nil, err
	}
	return &area, nil
}

// GetAreaList 获取区域列表，该方法会根据页码信息和排序信息返回分页后的区域信息。调用该方法需要传递的参数有4个，第一个为过滤信息，也就是要查询的区域信息或关键字；
// 第二个参数是分页信息；第三个参数是排序字段，第四个参数是是否倒序。如查询区域名称为“南京”，且要按照ID字段倒序排序，则参数 area.AreaName="南京", page.Page=1,
// page.PageInfo=10, order="id", desc=true
func (a *AreaService) GetAreaList(area *eaglesnest.Area, page request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&eaglesnest.Area{})
	var areas []eaglesnest.Area
	if area != nil {
		if area.AreaName != "" {
			db = db.Where("area_name LIKE ?", "%"+area.AreaName+"%")
		}
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
