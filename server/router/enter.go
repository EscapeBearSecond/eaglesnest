package router

import (
	"47.103.136.241/goprojects/curesan/server/router/curescan"
	"47.103.136.241/goprojects/curesan/server/router/example"
	"47.103.136.241/goprojects/curesan/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Curescan curescan.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
