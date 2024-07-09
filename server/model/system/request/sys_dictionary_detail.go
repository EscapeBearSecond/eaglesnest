package request

import (
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
