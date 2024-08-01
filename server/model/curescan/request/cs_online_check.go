package request

import "47.103.136.241/goprojects/curescan/server/model/common/request"

type SearchInfo struct {
	request.PageInfo
	TaskID   uint   `json:"taskId"`
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
