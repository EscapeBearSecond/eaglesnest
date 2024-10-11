package router

import (
	"github.com/EscapeBearSecond/curescan/server/router/curescan"
	"github.com/EscapeBearSecond/curescan/server/router/example"
	"github.com/EscapeBearSecond/curescan/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Curescan curescan.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
