package request

import (
	"47.103.136.241/goprojects/curescan/server/model/common/request"
	"47.103.136.241/goprojects/curescan/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
