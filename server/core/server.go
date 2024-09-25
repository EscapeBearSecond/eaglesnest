package core

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/initialize"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service/system"
	"fmt"
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

	fmt.Print(`
	欢迎使用 南京治煜漏扫平台
	当前版本:v1.0.0
	--------------------------------------版权声明--------------------------------------
	** 版权所有方：南京治煜开发团队 **
	** 版权持有公司：南京治煜信息科技有限公司 **
	** 版权声明：未经南京治煜信息科技有限公司书面许可，任何单位和个人不得以任何形式复制、传播、修改或商业使用本软件。 **
`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
