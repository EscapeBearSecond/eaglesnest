package request

import (
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
