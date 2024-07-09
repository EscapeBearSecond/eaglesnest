package system

import (
	"47.103.136.241/goprojects/curesan/server/config"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}
