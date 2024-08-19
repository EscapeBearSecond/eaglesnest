package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"47.103.136.241/goprojects/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type VulnRouter struct {
}

func (v *VulnRouter) InitVulnRouter(router *gin.RouterGroup) {
	vulnRouter := router.Group("vuln").Use(middleware.OperationRecord())
	vulnRouterApi := v1.ApiGroupApp.CurescanApiGroup.VulnApi
	statisticsApi := v1.ApiGroupApp.CurescanApiGroup.StatisticsApi
	{
		// vulnRouter.POST("", vulnRouterApi.Vuln)
		vulnRouter.POST("list", vulnRouterApi.GetVulnList)
		vulnRouter.GET("migrate", vulnRouterApi.MigrateTable)
		vulnRouter.GET("statistics", statisticsApi.GetVulnsInfo)
		vulnRouter.GET("common", statisticsApi.CommonVulnTopN)
	}
}
