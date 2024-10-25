package response

import (
	"github.com/EscapeBearSecond/eaglesnest/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
