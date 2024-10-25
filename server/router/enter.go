package router

import (
	"github.com/EscapeBearSecond/eaglesnest/server/router/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/router/example"
	"github.com/EscapeBearSecond/eaglesnest/server/router/system"
)

type RouterGroup struct {
	System     system.RouterGroup
	Example    example.RouterGroup
	Eaglesnest eaglesnest.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
