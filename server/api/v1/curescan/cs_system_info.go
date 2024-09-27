package curescan

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/response"
	"github.com/gin-gonic/gin"
)

type SystemInfoApi struct {
}

func (s *SystemInfoApi) GetSystemInfo(c *gin.Context) {
	info, err := systemInfoService.GetSystemInfo()
	if err != nil {
		response.FailWithMessage("获取系统信息失败", c)
		return
	}
	response.OkWithData(info, c)
}
