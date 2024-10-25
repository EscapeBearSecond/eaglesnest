package v1

import (
	"github.com/EscapeBearSecond/eaglesnest/server/api/v1/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/api/v1/example"
	"github.com/EscapeBearSecond/eaglesnest/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup     system.ApiGroup
	ExampleApiGroup    example.ApiGroup
	EaglesnestApiGroup eaglesnest.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
