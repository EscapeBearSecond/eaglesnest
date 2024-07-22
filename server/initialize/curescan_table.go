package initialize

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
)

func MigrateTables() {
	err := global.GVA_DB.AutoMigrate(&curescan.PortScan{}, &curescan.OnlineCheck{}, &curescan.JobResultItem{})
	if err != nil {
		panic(err)
	}
}
