package curescan

import (
	v1 "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/middleware"
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
