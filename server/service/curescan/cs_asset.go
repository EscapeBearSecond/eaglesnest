package curescan

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/request"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
)

type AssetService struct {
}

// BatchAdd 批量添加资产, 该方法适用于在资产扫描完成后, 已经得到了所有资产后, 调用该方法将资产全部添加到数据库中.
// 该方法分批次添加, 每个批次最多添加100条数据, 即如果有1000条待添加的资产, 该方法会分10次添加, 会产生10条SQL.
func (a *AssetService) BatchAdd(assets []*curescan.Asset) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "asset_ip"}, {Name: "created_by"}},
			DoUpdates: clause.AssignmentColumns([]string{"asset_name", "asset_area", "asset_type", "open_ports", "system_type", "ttl", "asset_model", "manufacturer"}),
		}).CreateInBatches(assets, 100).Error; err != nil {
			return err
		}
		return nil
	})
}
func (a *AssetService) BatchAddWithTransaction(tx *gorm.DB, assets []*curescan.Asset) error {
	// seen := make(map[string]struct{})
	// var result []*curescan.Asset
	// for _, asset := range assets {
	// 	key := fmt.Sprintf("%s%d", asset.AssetIP, asset.CreatedBy) // 根据唯一字段生成键
	// 	if _, ok := seen[key]; !ok {
	// 		seen[key] = struct{}{}
	// 		result = append(result, asset)
	// 	}
	// }
	return tx.Transaction(func(tx *gorm.DB) error {
		for _, asset := range assets {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "asset_ip"}, {Name: "created_by"}},
				DoUpdates: clause.AssignmentColumns([]string{"asset_name", "asset_area", "asset_type", "open_ports", "system_type", "ttl", "asset_model", "manufacturer", "area_name"}),
			}).Create(asset).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetAssetList 获取资产列表, 该方法会根据页码信息和排序信息返回分页后的资产信息. 调用该方法需要传递的参数有4个, 第一个为过滤信息, 也就是要查询的资产信息或关键字;
// 第二个参数是分页信息; 第三个参数是排序字段; 第四个参数是是否倒序. 如查询资产类型为"监控", 且要按照资产厂商字段倒序排序, 则参数 asset.AssetType="监控", page.Page=1,
// page.PageInfo=10, order="manufacturer", desc=true
func (a *AssetService) GetAssetList(asset *curescan.Asset, page request.PageInfo, order string, desc bool, allData bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Select("id", "asset_name", "asset_ip", "asset_area", "area_name", "asset_type", "open_ports", "system_type",
		"ttl", "asset_model", "manufacturer", "created_at", "updated_at", "deleted_at").Model(&curescan.Asset{})
	var assets []curescan.Asset
	if asset != nil {
		if asset.AssetName != "" {
			db = db.Where("asset_name LIKE ?", "%"+asset.AssetName+"%")
		}
		if asset.AreaName != "" {
			db = db.Where("area_name LIKE ?", "%"+asset.AreaName+"%")
		}
		if asset.AssetArea != 0 {
			db = db.Where("asset_area = ?", asset.AssetArea)
		}
		if asset.AssetIP != "" {
			db = db.Where("asset_ip LIKE ?", "%"+asset.AssetIP+"%")
		}
		if asset.Manufacturer != "" {
			db = db.Where("manufacturer LIKE ?", "%"+asset.Manufacturer+"%")
		}
		if asset.AssetModel != "" {
			db = db.Where("asset_model LIKE ?", "%"+asset.AssetModel+"%")
		}
		if asset.SystemType != "" {
			db = db.Where("system_type LIKE ?", "%"+asset.SystemType+"%")
		}
		if asset.AssetType != "" {
			db = db.Where("asset_type LIKE ?", "%"+asset.AssetType+"%")
		}
		if !allData {
			if asset.CreatedBy != 0 {
				db = db.Where("created_by = ?", asset.CreatedBy)
			}
		}
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
		orderMap["area_name"] = true
		orderMap["asset_area"] = true
		orderMap["manufacturer"] = true
		orderMap["system_type"] = true
		orderMap["asset_model"] = true
		orderMap["asset_ip"] = true
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
