package curescan

import (
	v1 "47.103.136.241/goprojects/curesan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type OnlineCheckRouter struct {
}

func (o *OnlineCheckRouter) InitOnlineCheckRouter(Router *gin.RouterGroup) {
	onlineCheckRouter := Router.Group("onlinecheck")

	onlineCheckRouterApi := v1.ApiGroupApp.CurescanApiGroup.OnlineCheckApi
	{
		onlineCheckRouter.GET("", onlineCheckRouterApi.GetInfoList)
	}
}
