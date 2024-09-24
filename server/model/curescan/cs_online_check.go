package curescan

import "47.103.136.241/goprojects/curescan/server/global"

type OnlineCheck struct {
	global.GvaModel
	global.CsModel
	IP      string `gorm:"column:ip;type:text;comment:主机" json:"ip"`
	System  string `gorm:"column:system;type:text;comment:系统" json:"system"`
	TTL     int    `gorm:"column:ttl;type:int8;comment:TTL" json:"ttl"`
	Active  bool   `gorm:"column:active;type:bool;comment:存活" json:"active"`
	TaskID  uint   `gorm:"column:task_id;type:int8;comment:任务ID" json:"taskId"`
	EntryID string `gorm:"column:entry_id;type:text;comment:entry ID" json:"entryId"`
}

func (OnlineCheck) TableName() string {
	return "cs_online_check"
}
