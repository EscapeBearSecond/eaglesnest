package curescan

import (
	v1 "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type SystemInfoRouter struct {
}

func (s *SystemInfoRouter) InitSystemInfoRouter(Router *gin.RouterGroup) {
	systemInfoRouter := Router.Group("systemInfo")
	systemInfoApi := v1.ApiGroupApp.CurescanApiGroup.SystemInfoApi
	{
		systemInfoRouter.GET("", systemInfoApi.GetSystemInfo)
	}
}
