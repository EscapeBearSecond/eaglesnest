package curescan

import (
	v1 "47.103.136.241/goprojects/curesan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type AreaRouter struct {
}

func (a *AreaRouter) InitAreaRouter(Router *gin.RouterGroup) {
	areaRouter := Router.Group("area")

	areaRouterApi := v1.ApiGroupApp.CurescanApiGroup.AreaApi
	{
		areaRouter.POST("", areaRouterApi.CreateArea) // 创建Area
		// areaRouter.POST("deleteArea", areaRouterApi.DeleteApi)               // 删除Area
		areaRouter.GET(":id", areaRouterApi.GetAreaById)       // 获取单条Area消息
		areaRouter.PUT("", areaRouterApi.UpdateArea)           // 更新Area
		areaRouter.DELETE(":id", areaRouterApi.DeleteAreaByID) // 删除选中Area

		// apiRouterWithoutRecord.POST("getAllAreas", areaRouterApi.GetAllApis) // 获取所有Area
		areaRouter.GET("", areaRouterApi.GetAreaList) // 获取Area列表

	}
}
