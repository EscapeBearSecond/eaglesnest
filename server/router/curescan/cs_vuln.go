package curescan

import (
	v1 "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type VulnRouter struct {
}

func (v *VulnRouter) InitVulnRouter(router *gin.RouterGroup) {
	vulnRouter := router.Group("vuln")
	vulnRouterApi := v1.ApiGroupApp.CurescanApiGroup.VulnApi
	statisticsApi := v1.ApiGroupApp.CurescanApiGroup.StatisticsApi
	{
		// vulnRouter.POST("", vulnRouterApi.Vuln)
		vulnRouter.POST("list", vulnRouterApi.GetVulnList)
		vulnRouter.GET("migrate", vulnRouterApi.MigrateTable)
		vulnRouter.GET("statistics", statisticsApi.GetVulnsInfo)
		vulnRouter.GET("common", statisticsApi.CommonVulnTopN)
		vulnRouter.POST("datamod", vulnRouterApi.DataMod)
	}
}
