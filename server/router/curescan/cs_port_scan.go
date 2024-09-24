package curescan

import (
	v1 "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type PortScanRouter struct {
}

func (p *PortScanRouter) InitPortScanRouter(Router *gin.RouterGroup) {
	portScanRouter := Router.Group("portscan").Use(middleware.OperationRecord())
	portScanRouterApi := v1.ApiGroupApp.CurescanApiGroup.PortScanApi
	{
		portScanRouter.GET("", portScanRouterApi.GetInfoList)
	}
}
