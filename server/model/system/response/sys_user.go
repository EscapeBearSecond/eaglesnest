package response

import (
	"47.103.136.241/goprojects/curescan/server/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}
