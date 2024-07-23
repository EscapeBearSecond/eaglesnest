package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"fmt"
)

type JobResultService struct {
}

func (j *JobResultService) BatchAdd(data []*curescan.JobResultItem) error {
	fmt.Println("待插入数据: ", len(data))
	tx := global.GVA_DB.Model(&curescan.JobResultItem{}).CreateInBatches(data, 100)
	fmt.Println("影响行数: ", tx.RowsAffected)
	return tx.Error
}
