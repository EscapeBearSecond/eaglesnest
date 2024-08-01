package system

import (
	"47.103.136.241/goprojects/curescan/server/global"
)

type JwtBlacklist struct {
	global.GvaModel
	Jwt string `gorm:"type:text;comment:jwt"`
}
