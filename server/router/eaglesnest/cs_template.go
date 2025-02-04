package eaglesnest

import (
	v1 "github.com/EscapeBearSecond/eaglesnest/server/api/v1"
	"github.com/EscapeBearSecond/eaglesnest/server/middleware"
	"github.com/gin-gonic/gin"
)

type TemplateRouter struct {
}

func (t *TemplateRouter) InitTemplateRouter(Router *gin.RouterGroup) {
	templateRouter := Router.Group("template")
	templateRouterApi := v1.ApiGroupApp.EaglesnestApiGroup.TemplateApi

	{
		templateRouter.GET("", templateRouterApi.MigrateTable)
		templateRouter.POST("", templateRouterApi.CreateTemplate).Use(middleware.OperationRecord()) // 创建Template

		templateRouter.GET(":id", templateRouterApi.GetTemplateById)                                     // 获取单条Template消息
		templateRouter.PUT("", templateRouterApi.UpdateTemplate).Use(middleware.OperationRecord())       // 更新Template
		templateRouter.DELETE(":id", templateRouterApi.DeleteTemplate).Use(middleware.OperationRecord()) // 删除选中Template
		templateRouter.POST("importTemplateContent", templateRouterApi.ImportTemplateContent).Use(middleware.OperationRecord())

		// apiRouterWithoutRecord.POST("getAllTemplates", TemplateRouterApi.GetAllApis) // 获取所有Template
		templateRouter.POST("list", templateRouterApi.GetTemplateList) // 获取Template列表
		templateRouter.POST("imports", templateRouterApi.ImportTemplates)
		templateRouter.GET("tags", templateRouterApi.TemplateTags)
		templateRouter.GET("lll", templateRouterApi.LLL)
		templateRouter.POST("update", templateRouterApi.UploadFromZip)

	}
}
