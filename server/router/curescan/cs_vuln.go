package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type VulnRouter struct {
}

func (v *VulnRouter) InitVulnRouter(router *gin.RouterGroup) {
	vulnRouter := router.Group("vuln")
	vulnRouterApi := v1.ApiGroupApp.CurescanApiGroup.VulnApi
	{
		// vulnRouter.POST("", vulnRouterApi.Vuln)
		vulnRouter.POST("list", vulnRouterApi.GetVulnList)
		vulnRouter.GET("migrate", vulnRouterApi.MigrateTable)
	}
}
