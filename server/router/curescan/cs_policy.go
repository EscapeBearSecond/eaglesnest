package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type PolicyRouter struct {
}

func (p *PolicyRouter) InitPolicyRouter(Router *gin.RouterGroup) {
	policyRouter := Router.Group("policy")
	policyRouterApi := v1.ApiGroupApp.CurescanApiGroup.PolicyApi
	{
		// policyRouter.GET("", policyRouterApi.MigrateTable)
		policyRouter.POST("", policyRouterApi.CreatePolicy)
		policyRouter.PUT("", policyRouterApi.UpdatePolicy)
		policyRouter.DELETE(":id", policyRouterApi.DeletePolicy)
		policyRouter.POST("list", policyRouterApi.GetPolicyList)
		policyRouter.GET(":id", policyRouterApi.GetPolicyById)
	}
}
