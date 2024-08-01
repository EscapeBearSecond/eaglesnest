package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"github.com/lib/pq"
)

type PortScan struct {
	global.GvaModel
	global.CsModel
	TaskID  uint          `gorm:"column:task_id;type:int8;comment:任务ID" json:"taskId"`
	IP      string        `gorm:"column:ip;type:varchar(20);comment:主机" json:"ip"`
	Ports   pq.Int64Array `gorm:"column:port;type:int8[];comment:端口" json:"ports"`
	EntryID string        `gorm:"column:entry_id;type:varchar(255);comment:entry ID" json:"entryId"`
}

func (PortScan) TableName() string {
	return "cs_port_scan"
}
