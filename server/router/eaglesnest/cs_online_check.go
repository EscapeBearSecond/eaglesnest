package eaglesnest

import (
	v1 "github.com/EscapeBearSecond/eaglesnest/server/api/v1"
	"github.com/EscapeBearSecond/eaglesnest/server/middleware"
	"github.com/gin-gonic/gin"
)

type OnlineCheckRouter struct {
}

func (o *OnlineCheckRouter) InitOnlineCheckRouter(Router *gin.RouterGroup) {
	onlineCheckRouter := Router.Group("onlinecheck").Use(middleware.OperationRecord())

	onlineCheckRouterApi := v1.ApiGroupApp.EaglesnestApiGroup.OnlineCheckApi
	{
		onlineCheckRouter.GET("", onlineCheckRouterApi.GetInfoList)
	}
}
