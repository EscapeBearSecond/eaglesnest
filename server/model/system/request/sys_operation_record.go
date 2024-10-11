package request

import (
	"github.com/EscapeBearSecond/curescan/server/model/common/request"
	"github.com/EscapeBearSecond/curescan/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
