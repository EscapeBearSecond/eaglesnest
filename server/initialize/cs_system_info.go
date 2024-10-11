package initialize

import (
	"github.com/EscapeBearSecond/curescan/server/global"
	"github.com/EscapeBearSecond/curescan/server/service/curescan"
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
