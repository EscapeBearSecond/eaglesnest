package curescan

import "47.103.136.241/goprojects/curesan/server/global"

type OnlineCheck struct {
	global.GvaModel
	IP     string `gorm:"column:ip;type:varchar(20);comment:主机" json:"ip"`
	System string `gorm:"column:system;type:varchar(20);comment:系统" json:"system"`
	TTL    string `gorm:"column:ttl;type:varchar(20);comment:TTL" json:"ttl"`
	Active string `gorm:"column:active;type:char(1);comment:存活" json:"active"`
	TaskID uint   `gorm:"column:task_id;type:int8;comment:任务ID" json:"taskId"`
}

func (OnlineCheck) TableName() string {
	return "cs_online_check"
}
