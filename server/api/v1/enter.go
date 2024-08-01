package v1

import (
	"47.103.136.241/goprojects/curescan/server/api/v1/curescan"
	"47.103.136.241/goprojects/curescan/server/api/v1/example"
	"47.103.136.241/goprojects/curescan/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	CurescanApiGroup curescan.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
