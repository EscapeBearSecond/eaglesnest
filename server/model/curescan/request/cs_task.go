package request

import (
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
)

type CreateTask struct {
	TaskName   string   `json:"taskName"`
	TaskDesc   string   `json:"taskDesc"`
	Status     int      `json:"status"`
	TargetIP   []string `json:"targetIp"`
	PolicyID   uint     `json:"policyId"`
	TaskPlan   int      `json:"taskPlan"`
	PlanConfig string   `json:"planConfig"`
}

type PlanConfig struct {
	Date      string `json:"date"`
	Time      string `json:"time"`
	Frequency int    `json:"frequency"` // 1:每天 2:每周 3:每月
}

type UpdateTask struct {
	ID uint `json:"id"`
	CreateTask
}

type SearchTask struct {
	curescan.Task
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
