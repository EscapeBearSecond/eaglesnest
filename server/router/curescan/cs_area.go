package curescan

import (
	v1 "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type AreaRouter struct {
}

func (a *AreaRouter) InitAreaRouter(Router *gin.RouterGroup) {
	areaRouter := Router.Group("area")

	areaRouterApi := v1.ApiGroupApp.CurescanApiGroup.AreaApi
	{
		areaRouter.POST("", areaRouterApi.CreateArea).Use(middleware.OperationRecord()) // 创建Area
		// areaRouter.POST("deleteArea", areaRouterApi.DeleteApi)               // 删除Area
		areaRouter.GET(":id", areaRouterApi.GetAreaById)                                         // 获取单条Area消息
		areaRouter.PUT("", areaRouterApi.UpdateArea).Use(middleware.OperationRecord())           // 更新Area
		areaRouter.DELETE(":id", areaRouterApi.DeleteAreaByID).Use(middleware.OperationRecord()) // 删除选中Area

		// apiRouterWithoutRecord.POST("getAllAreas", areaRouterApi.GetAllApis) // 获取所有Area
		areaRouter.POST("list", areaRouterApi.GetAreaList) // 获取Area列表

	}
}
