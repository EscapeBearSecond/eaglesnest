package initialize

import (
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/service/eaglesnest"
	"go.uber.org/zap"
)

var systemInfoService = &eaglesnest.SystemInfoService{}

func InitSystemInfo() {
	if global.GVA_DB != nil {
		err := systemInfoService.InitSystemInfo()
		if err != nil {
			global.GVA_LOG.Error("初始化系统信息失败", zap.Error(err))
		}
	}
}
