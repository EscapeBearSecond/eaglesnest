package core

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/initialize"
	"47.103.136.241/goprojects/curesan/server/service/system"
	"fmt"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}

	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)

	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 南京治煜漏扫平台
	当前版本:v1.0.0
	--------------------------------------版权声明--------------------------------------
	** 版权所有方：南京治煜开发团队 **
	** 版权持有公司：南京治煜信息科技有限公司 **
`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
