package system

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/config"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}
