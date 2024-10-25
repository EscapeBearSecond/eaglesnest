package eaglesnest

import (
	v1 "github.com/EscapeBearSecond/eaglesnest/server/api/v1"
	"github.com/gin-gonic/gin"
)

type VulnRouter struct {
}

func (v *VulnRouter) InitVulnRouter(router *gin.RouterGroup) {
	vulnRouter := router.Group("vuln")
	vulnRouterApi := v1.ApiGroupApp.EaglesnestApiGroup.VulnApi
	statisticsApi := v1.ApiGroupApp.EaglesnestApiGroup.StatisticsApi
	{
		// vulnRouter.POST("", vulnRouterApi.Vuln)
		vulnRouter.POST("list", vulnRouterApi.GetVulnList)
		vulnRouter.GET("migrate", vulnRouterApi.MigrateTable)
		vulnRouter.GET("statistics", statisticsApi.GetVulnsInfo)
		vulnRouter.GET("common", statisticsApi.CommonVulnTopN)
		vulnRouter.POST("datamod", vulnRouterApi.DataMod)
	}
}
