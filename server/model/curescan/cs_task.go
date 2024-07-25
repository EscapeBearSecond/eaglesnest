package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"github.com/lib/pq"
)

type Task struct {
	global.GvaModel
	TaskName   string         `gorm:"column:task_name;type:varchar(50);not null;uniqueIndex;comment:任务名称"`
	TaskDesc   string         `gorm:"column:task_desc;type:varchar(100);comment:任务描述"`
	Status     int            `gorm:"column:status;type:int2;comment:执行状态"` // 执行状态 0 创建、1 执行中、2 执行完成、3 执行失败
	TargetIP   pq.StringArray `gorm:"column:target_ip;type:text[];comment:目标IP"`
	PolicyID   uint           `gorm:"column:policy_id;type:int8;comment:策略ID"`
	TaskPlan   int            `gorm:"column:task_plan;type:int2;comment:任务计划"` // 任务计划 0 立即执行、1 稍后执行、 2定时执行
	PlanConfig string         `gorm:"column:plan_config;type:text;comment:计划配置"`
	Executions uint           `gorm:"column:executions;type:int8;comment:执行次数"`
}

func (Task) TableName() string {
	return "cs_task"
}
