package router

import (
	"47.103.136.241/goprojects/curescan/server/router/curescan"
	"47.103.136.241/goprojects/curescan/server/router/example"
	"47.103.136.241/goprojects/curescan/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Curescan curescan.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
