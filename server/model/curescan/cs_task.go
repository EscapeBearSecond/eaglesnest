package curescan

import (
	"github.com/EscapeBearSecond/curescan/server/global"
	"github.com/lib/pq"
)

type Task struct {
	global.GvaModel
	global.CsModel
	TaskName    string         `json:"taskName" gorm:"column:task_name;type:text;not null;comment:任务名称"`
	TaskDesc    string         `json:"taskDesc" gorm:"column:task_desc;type:text;comment:任务描述"`
	Status      int            `json:"status" gorm:"column:status;type:int2;comment:执行状态"` // 执行状态 0 创建、1 执行中、2 执行完成、3 执行失败
	TargetIP    pq.StringArray `json:"targetIp" gorm:"column:target_ip;type:text[];comment:目标IP"`
	PolicyID    uint           `json:"policyId" gorm:"column:policy_id;type:int8;comment:策略ID"`
	PolicyName  string         `json:"policyName" gorm:"-"`
	TaskPlan    int            `json:"taskPlan" gorm:"column:task_plan;type:int2;comment:任务计划"` // 任务计划 1 立即执行、2 稍后执行、 3定时执行
	PlanConfig  string         `json:"planConfig" gorm:"column:plan_config;type:text;comment:计划配置"`
	Executions  uint           `json:"executions" gorm:"column:executions;type:int8;comment:执行次数"`
	EntryID     string         `json:"entryId" gorm:"column:entry_id;type:text;comment:entry ID"`
	Flag        string         `json:"flag" gorm:"column:flag;type:text;comment:flag"` // 用来标记是通过区域创建还是自定义ip创建
	AreaIDArray pq.Int64Array  `json:"areaIdArray" gorm:"column:area_id_array;type:int8[];comment:区域ID"`
	StartAt     string         `json:"startAt"`
	EndAt       string         `json:"endAt"`
}

func (Task) TableName() string {
	return "cs_task"
}
