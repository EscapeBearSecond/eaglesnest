package global

import "errors"

var (
	HasExisted     = errors.New("already exists")
	NotAllowed     = errors.New("not allowed")
	NoPermission   = errors.New("no permission")
	NoDataFound    = errors.New("no data found")
	ParamError     = errors.New("parameter error")
	NoDataSelected = errors.New("no data selected")
)
