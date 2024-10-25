package response

import "github.com/EscapeBearSecond/eaglesnest/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
