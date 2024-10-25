package eaglesnest

import (
	v1 "github.com/EscapeBearSecond/eaglesnest/server/api/v1"
	"github.com/gin-gonic/gin"
)

type AssetRouter struct {
}

func (a *AssetRouter) InitAssetRouter(Router *gin.RouterGroup) {
	assetRouter := Router.Group("asset")
	assetRouterApi := v1.ApiGroupApp.EaglesnestApiGroup.AssetApi
	statisticsApi := v1.ApiGroupApp.EaglesnestApiGroup.StatisticsApi
	{
		assetRouter.GET("", assetRouterApi.MigrateTable)
		assetRouter.POST("/batchAdd", assetRouterApi.BatchAdd)
		assetRouter.POST("/list", assetRouterApi.GetAssetList)
		assetRouter.GET("highrisk", statisticsApi.AssetTopN)
	}
}
