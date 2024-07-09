package response

import "47.103.136.241/goprojects/curesan/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
