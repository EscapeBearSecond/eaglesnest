package response

import (
	"47.103.136.241/goprojects/curescan/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
