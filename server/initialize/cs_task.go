package initialize

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"47.103.136.241/goprojects/curescan/server/service/curescan"
	"47.103.136.241/goprojects/curescan/server/utils"
	"context"
	"go.uber.org/zap"
	"strconv"
	"sync"
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
		// 从队列中取出但是还没来得及执行的task，任务状态没有更新
		var ids []int
		err = global.GVA_DB.Exec("select id from cs_task where status = ?", common.Waiting).Scan(&ids).Error
		if err != nil {
			global.GVA_LOG.Error("RecoverTask get task id failed", zap.Error(err))
			return
		}
		for _, id := range ids {
			utils.RemoveValueFromList(global.GVA_REDIS, "taskQueue", strconv.Itoa(id))
		}

	}
}

var taskService = &curescan.TaskService{}

func ExecuteTask() {
	go func() {
		var wg sync.WaitGroup
		for {
			ids, err := global.GVA_REDIS.BLPop(context.Background(), 0, "taskQueue").Result()
			if err != nil {
				continue
			}
			taskId, _ := strconv.Atoi(ids[1])
			wg.Add(1)
			taskService.ExecuteTask(taskId, &wg)
			wg.Wait()
		}
	}()
}
