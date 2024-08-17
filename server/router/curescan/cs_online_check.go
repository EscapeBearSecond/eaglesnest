package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"47.103.136.241/goprojects/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type OnlineCheckRouter struct {
}

func (o *OnlineCheckRouter) InitOnlineCheckRouter(Router *gin.RouterGroup) {
	onlineCheckRouter := Router.Group("onlinecheck").Use(middleware.OperationRecord())

	onlineCheckRouterApi := v1.ApiGroupApp.CurescanApiGroup.OnlineCheckApi
	{
		onlineCheckRouter.GET("", onlineCheckRouterApi.GetInfoList)
	}
}
