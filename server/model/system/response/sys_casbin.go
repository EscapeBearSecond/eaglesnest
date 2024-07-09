package response

import (
	"47.103.136.241/goprojects/curesan/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
