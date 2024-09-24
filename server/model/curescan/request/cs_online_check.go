package request

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/request"

type SearchInfo struct {
	request.PageInfo
	TaskID   uint   `json:"taskId"`
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
