package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"47.103.136.241/goprojects/curescan/server/middleware"
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
