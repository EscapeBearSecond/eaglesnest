package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"47.103.136.241/goprojects/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type AssetRouter struct {
}

func (a *AssetRouter) InitAssetRouter(Router *gin.RouterGroup) {
	assetRouter := Router.Group("asset").Use(middleware.OperationRecord())
	assetRouterApi := v1.ApiGroupApp.CurescanApiGroup.AssetApi

	{
		assetRouter.GET("", assetRouterApi.MigrateTable)
		assetRouter.POST("/batchAdd", assetRouterApi.BatchAdd)
		assetRouter.POST("/list", assetRouterApi.GetAssetList)
	}
}
