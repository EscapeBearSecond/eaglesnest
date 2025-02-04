package initialize

import (
	"os"

	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/model/example"
	"github.com/EscapeBearSecond/eaglesnest/server/model/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "pgsql":
		return GormPgSql()
	default:
		return GormPgSql()
	}
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},
		system.SysExportTemplate{},
		system.Condition{},
		system.JoinTemplate{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},

		eaglesnest.Area{},
		eaglesnest.Asset{},
		eaglesnest.JobResultItem{},
		eaglesnest.OnlineCheck{},
		eaglesnest.Policy{},
		eaglesnest.PortScan{},
		eaglesnest.Task{},
		eaglesnest.Template{},
		eaglesnest.Vuln{},
		eaglesnest.SystemInfo{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel(db)

	if err != nil {
		global.GVA_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
