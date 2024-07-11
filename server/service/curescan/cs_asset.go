package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"fmt"
)

type AssetService struct {
}

func (a *AssetService) BatchAdd(assets []*curescan.Asset) error {
	return global.GVA_DB.Model(&curescan.Asset{}).CreateInBatches(assets, 100).Error
}

func (a *AssetService) GetAreaList(asset curescan.Asset, page request.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.Asset{})
	var assets []curescan.Asset
	if asset.AssetName != "" {
		db = db.Where("asset_name LIKE ?", "%"+asset.AssetName+"%")
	}
	if asset.AssetArea != 0 {
		db = db.Where("asset_area = ?", asset.AssetArea)
	}
	if asset.Manufacturer != "" {
		db = db.Where("manufacturer LIKE ?", "%"+asset.Manufacturer+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return assets, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool)
		orderMap["id"] = true
		orderMap["asset_name"] = true
		orderMap["asset_area"] = true
		orderMap["manufacturer"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return assets, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&assets).Error
	return assets, total, err
}
