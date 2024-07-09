package router

import (
	"47.103.136.241/goprojects/gin-vue-admin/server/router/example"
	"47.103.136.241/goprojects/gin-vue-admin/server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
