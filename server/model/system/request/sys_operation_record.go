package request

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/request"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
