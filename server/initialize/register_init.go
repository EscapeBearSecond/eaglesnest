package initialize

import (
	_ "github.com/EscapeBearSecond/eaglesnest/server/source/eaglesnest"
	_ "github.com/EscapeBearSecond/eaglesnest/server/source/example"
	_ "github.com/EscapeBearSecond/eaglesnest/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
