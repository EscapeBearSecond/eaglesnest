package curescan

import (
	v1 "47.103.136.241/goprojects/curesan/server/api/v1"
	"47.103.136.241/goprojects/curesan/server/middleware"
	"github.com/gin-gonic/gin"
)

type AreaRouter struct {
}

func (a *AreaRouter) InitAreaRouter(Router *gin.RouterGroup) {
	areaRouter := Router.Group("area").Use(middleware.OperationRecord())

	areaRouterApi := v1.ApiGroupApp.CurescanApiGroup.AreaApi
	{
		areaRouter.POST("", areaRouterApi.CreateArea) // 创建Area
		// areaRouter.POST("deleteArea", areaRouterApi.DeleteApi)               // 删除Area
		// areaRouter.POST("getAreaById", areaRouterApi.GetApiById)             // 获取单条Area消息
		// areaRouter.POST("updateArea", areaRouterApi.UpdateApi)               // 更新Area
		// areaRouter.DELETE("deleteAreasByIds", areaRouterApi.DeleteApisByIds) // 删除选中Area
	}
	{
		// apiRouterWithoutRecord.POST("getAllAreas", areaRouterApi.GetAllApis) // 获取所有Area
		// apiRouterWithoutRecord.POST("getAreaList", areaRouterApi.GetApiList) // 获取Area列表
	}

}
