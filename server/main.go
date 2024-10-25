package main

import (
	"github.com/EscapeBearSecond/eaglesnest/server/model/system/request"
	"github.com/EscapeBearSecond/eaglesnest/server/service/system"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
	_ "net/http/pprof"

	"github.com/EscapeBearSecond/eaglesnest/server/core"
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       eaglesnest Swagger API接口文档
// @version                     v2.6.5
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	initialize.FalconEngine()

	if global.GVA_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	} else {
		dbInfo := request.InitDB{
			AdminPassword: "@123456qwer",
			DBType:        global.GVA_CONFIG.System.DbType,
			DBName:        global.GVA_CONFIG.Pgsql.Dbname,
			Host:          global.GVA_CONFIG.Pgsql.Path,
			Port:          global.GVA_CONFIG.Pgsql.Port,
			UserName:      global.GVA_CONFIG.Pgsql.Username,
			Password:      global.GVA_CONFIG.Pgsql.Password,
			DBPath:        global.GVA_CONFIG.Pgsql.Path,
		}
		global.GVA_LOG.Info("数据库信息", zap.Any("dbInfo", dbInfo))
		service := system.InitDBService{}
		err := service.InitDB(dbInfo)
		if err != nil {
			global.GVA_LOG.Error("初始化数据库失败", zap.Any("err", err))
		}
	}
	core.RunWindowsServer()
}
