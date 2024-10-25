package request

import (
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
