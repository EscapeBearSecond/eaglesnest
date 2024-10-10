package global

import "errors"

var (
	HasExisted     = errors.New("数据已存在")
	NotAllowed     = errors.New("操作不被允许")
	NoPermission   = errors.New("无权限")
	NoDataFound    = errors.New("数据不存在")
	ParamError     = errors.New("参数错误")
	NoDataSelected = errors.New("未选择数据")
	PasswordError  = errors.New("密码错误")
)
