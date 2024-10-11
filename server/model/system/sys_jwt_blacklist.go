package system

import (
	"github.com/EscapeBearSecond/curescan/server/global"
)

type JwtBlacklist struct {
	global.GvaModel
	Jwt string `gorm:"type:text;comment:jwt"`
}
