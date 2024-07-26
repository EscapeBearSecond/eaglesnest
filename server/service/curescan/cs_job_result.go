package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"gorm.io/gorm"
)

type JobResultService struct {
}

func (j *JobResultService) BatchAdd(data []*curescan.JobResultItem) error {
	return global.GVA_DB.Model(&curescan.JobResultItem{}).CreateInBatches(data, 100).Error
}

func (j *JobResultService) BatchAddWithTransaction(tx *gorm.DB, data []*curescan.JobResultItem) error {
	return tx.Model(&curescan.JobResultItem{}).CreateInBatches(data, 100).Error
}
