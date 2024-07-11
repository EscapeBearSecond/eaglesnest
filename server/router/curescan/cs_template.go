package curescan

import (
	v1 "47.103.136.241/goprojects/curesan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type TemplateRouter struct {
}

func (t *TemplateRouter) InitTemplateRouter(Router *gin.RouterGroup) {
	templateRouter := Router.Group("template")
	templateRouterApi := v1.ApiGroupApp.CurescanApiGroup.TemplateApi

	{
		// templateRouter.GET("", templateRouterApi.MigrateTable)
		templateRouter.POST("", templateRouterApi.CreateTemplate) // 创建Template

		templateRouter.GET(":id", templateRouterApi.GetTemplateById)   // 获取单条Area消息
		templateRouter.PUT("", templateRouterApi.UpdateTemplate)       // 更新Area
		templateRouter.DELETE(":id", templateRouterApi.DeleteTemplate) // 删除选中Template
		templateRouter.POST("/importTemplateContent", templateRouterApi.ImportTemplateContent)

		// apiRouterWithoutRecord.POST("getAllAreas", areaRouterApi.GetAllApis) // 获取所有Area
		templateRouter.GET("", templateRouterApi.GetTemplateList) // 获取Area列表

	}
}
