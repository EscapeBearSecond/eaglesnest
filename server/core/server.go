package core

import (
	"fmt"
	"github.com/EscapeBearSecond/eaglesnest/server/core/meta"
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/initialize"
	"github.com/EscapeBearSecond/eaglesnest/server/service/system"
	"go.uber.org/zap"
	_ "net/http/pprof"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}
	initialize.RecoverTask()
	initialize.ExecuteTask()

	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)

	global.GVA_LOG.Info("server run success on ", zap.String("address", address))
	fmt.Print("BuildVer:", meta.BuildVer)
	fmt.Printf(`
	欢迎使用 DYG漏扫平台
	当前版本:v%s
`, meta.BuildVer)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
