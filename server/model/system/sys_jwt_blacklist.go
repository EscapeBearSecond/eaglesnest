package system

import (
	"github.com/EscapeBearSecond/eaglesnest/server/global"
)

type JwtBlacklist struct {
	global.GvaModel
	Jwt string `gorm:"type:text;comment:jwt"`
}
