package response

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
