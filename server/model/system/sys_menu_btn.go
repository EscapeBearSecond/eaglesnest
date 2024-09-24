package system

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"

type SysBaseMenuBtn struct {
	global.GvaModel
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}
