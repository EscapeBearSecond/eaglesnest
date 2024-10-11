package response

import "github.com/EscapeBearSecond/curescan/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
