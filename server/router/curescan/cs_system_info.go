package curescan

import (
	v1 "github.com/EscapeBearSecond/curescan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type SystemInfoRouter struct {
}

func (s *SystemInfoRouter) InitSystemInfoRouter(Router *gin.RouterGroup) {
	systemInfoRouter := Router.Group("systemInfo")
	systemInfoApi := v1.ApiGroupApp.CurescanApiGroup.SystemInfoApi
	{
		systemInfoRouter.GET("", systemInfoApi.GetSystemInfo)
		systemInfoRouter.POST("/license", systemInfoApi.UpdateLicense)
	}
}
