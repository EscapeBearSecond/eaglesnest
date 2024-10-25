package eaglesnest

import (
	v1 "github.com/EscapeBearSecond/eaglesnest/server/api/v1"
	"github.com/EscapeBearSecond/eaglesnest/server/middleware"
	"github.com/gin-gonic/gin"
)

type PortScanRouter struct {
}

func (p *PortScanRouter) InitPortScanRouter(Router *gin.RouterGroup) {
	portScanRouter := Router.Group("portscan").Use(middleware.OperationRecord())
	portScanRouterApi := v1.ApiGroupApp.EaglesnestApiGroup.PortScanApi
	{
		portScanRouter.GET("", portScanRouterApi.GetInfoList)
	}
}
