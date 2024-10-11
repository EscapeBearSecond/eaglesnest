package response

import (
	"github.com/EscapeBearSecond/curescan/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
