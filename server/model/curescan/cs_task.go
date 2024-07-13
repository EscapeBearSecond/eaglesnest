package curescan

import "47.103.136.241/goprojects/curesan/server/global"

type Task struct {
	global.GvaModel
	Name string `gorm:"type:varchar(255);not null;comment:任务名称"`
	Desc string `gorm:"type:varchar(255);comment:任务描述"`
	Status int `gorm:"type:tinyint;default:0;comment:任务状态"`
	Result string `gorm:"type:varchar(255);comment:任务结果"`
}