package v1

import (
	"github.com/EscapeBearSecond/curescan/server/api/v1/curescan"
	"github.com/EscapeBearSecond/curescan/server/api/v1/example"
	"github.com/EscapeBearSecond/curescan/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	CurescanApiGroup curescan.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
