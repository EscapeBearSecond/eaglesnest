package eaglesnest

import (
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/response"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/request"
	"github.com/EscapeBearSecond/eaglesnest/server/service/system"
	"github.com/EscapeBearSecond/eaglesnest/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AssetApi struct {
}

func (a *AssetApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&eaglesnest.Asset{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// BatchAdd
// 批量添加资产, 适用于资产发现后调用该接口
func (a *AssetApi) BatchAdd(c *gin.Context) {
	var assets []*eaglesnest.Asset
	err := c.ShouldBindJSON(&assets)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = assetService.BatchAdd(assets)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (a *AssetApi) GetAssetList(c *gin.Context) {
	var searchAsset request.SearchAsset
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("查看资产列表失败", zap.Error(err))
		}
	}()
	err = c.ShouldBindJSON(&searchAsset)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(searchAsset.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if searchAsset.Asset == nil {
		searchAsset.Asset = &eaglesnest.Asset{}
	}
	searchAsset.CreatedBy = utils.GetUserID(c)
	searchAsset.Asset.CreatedBy = utils.GetUserID(c)
	allData := system.HasAllDataAuthority(c)
	list, total, err := assetService.GetAssetList(searchAsset.Asset, searchAsset.PageInfo, searchAsset.OrderKey, searchAsset.Desc, allData)
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchAsset.Page,
		PageSize: searchAsset.PageSize,
	}, "获取成功", c)
}
