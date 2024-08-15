package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"fmt"
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

func (j *JobResultService) GetJobResultList(info *request.SearchJobResult) (list interface{}, total int64, err error) {
	jobResult := info.JobResultItem
	page := info.PageInfo
	order := info.OrderKey
	desc := info.Desc
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.JobResultItem{})
	var jobResultList []*curescan.JobResultItem
	if jobResult.Kind != "" {
		db = db.Where("kind = ?", jobResult.Kind)
	}
	if jobResult.Name != "" {
		db = db.Where("name = ?", jobResult.Name)
	}
	if jobResult.Severity != "" {
		db = db.Where("severity = ?", jobResult.Severity)
	}
	err = db.Count(&total).Error
	if err != nil {
		return jobResultList, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool)
		orderMap["id"] = true
		orderMap["kind"] = true
		orderMap["name"] = true
		orderMap["severity"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return jobResultList, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&jobResultList).Error
	return jobResultList, total, err
}
