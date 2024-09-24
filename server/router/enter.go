package router

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/router/curescan"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/router/example"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Curescan curescan.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
