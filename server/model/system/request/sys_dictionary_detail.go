package request

import (
	"47.103.136.241/goprojects/curescan/server/model/common/request"
	"47.103.136.241/goprojects/curescan/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
