package request

import (
	"47.103.136.241/goprojects/curescan/server/model/common/request"
)

type CreateTask struct {
	TaskName   string   `json:"taskName"`
	TaskDesc   string   `json:"taskDesc"`
	Status     int      `json:"status"`
	TargetIP   []string `json:"targetIp"`
	Flag       string   `json:"flag"`
	PolicyID   uint     `json:"policyId"`
	TaskPlan   int      `json:"taskPlan,string"`
	PlanConfig string   `json:"planConfig"`
}

type UpdateTask struct {
	ID uint `json:"ID"`
	CreateTask
}

type SearchTask struct {
	TaskName string `json:"taskName"`
	Status   int    `json:"status"`
	TaskPlan []int  `json:"taskPlan"`
	PolicyId int    `json:"policyId"`
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
