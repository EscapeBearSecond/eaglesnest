package response

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/system"

type SysAPIResponse struct {
	Api system.SysApi `json:"api"`
}

type SysAPIListResponse struct {
	Apis []system.SysApi `json:"apis"`
}
