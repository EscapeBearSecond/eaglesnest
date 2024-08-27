package initialize

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"go.uber.org/zap"
)

// RecoverTask 将因系统出现崩溃或异常而导致的任务恢复 将执行中改为失败
func RecoverTask() {
	if global.GVA_DB != nil {
		err := global.GVA_DB.Exec("update cs_task set status = ? where status = ?", common.Failed, common.Running).Error
		if err != nil {
			global.GVA_LOG.Error("RecoverTask common task failed", zap.Error(err))
		}
		err = global.GVA_DB.Exec("update cs_task set status = ? where status = ?", common.TimeStopped, common.TimeRunning).Error
		if err != nil {
			global.GVA_LOG.Error("RecoverTask timing task failed", zap.Error(err))
		}
	}
}
