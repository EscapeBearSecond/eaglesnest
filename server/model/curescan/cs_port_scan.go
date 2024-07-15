package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"github.com/lib/pq"
)

type PortScan struct {
	global.GvaModel
	TaskID uint           `gorm:"column:task_id;type:int8;comment:任务ID" json:"taskId"`
	IP     string         `gorm:"column:ip;type:varchar(20);comment:主机" json:"ip"`
	Ports  pq.StringArray `gorm:"column:port;type:int8[];comment:端口" json:"ports"`
}
