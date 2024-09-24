package curescan

import (
	v1 "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type PolicyRouter struct {
}

func (p *PolicyRouter) InitPolicyRouter(Router *gin.RouterGroup) {
	policyRouter := Router.Group("policy")
	policyRouterApi := v1.ApiGroupApp.CurescanApiGroup.PolicyApi
	{
		// policyRouter.GET("", policyRouterApi.MigrateTable)
		policyRouter.POST("", policyRouterApi.CreatePolicy).Use(middleware.OperationRecord())
		policyRouter.PUT("", policyRouterApi.UpdatePolicy).Use(middleware.OperationRecord())
		policyRouter.DELETE(":id", policyRouterApi.DeletePolicy).Use(middleware.OperationRecord())
		policyRouter.POST("list", policyRouterApi.GetPolicyList)
		policyRouter.GET(":id", policyRouterApi.GetPolicyById)
	}
}
