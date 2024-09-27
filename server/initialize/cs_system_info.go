package initialize

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service/curescan"
	"go.uber.org/zap"
)

var systemInfoService = &curescan.SystemInfoService{}

func InitSystemInfo() {
	if global.GVA_DB != nil {
		err := systemInfoService.InitSystemInfo()
		if err != nil {
			global.GVA_LOG.Error("初始化系统信息失败", zap.Error(err))
		}
	}
}
