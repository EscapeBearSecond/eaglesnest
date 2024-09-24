package response

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
