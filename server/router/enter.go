package router

import (
	"47.103.136.241/goprojects/curesan/server/router/example"
	"47.103.136.241/goprojects/curesan/server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
