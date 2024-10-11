package initialize

import (
	_ "github.com/EscapeBearSecond/curescan/server/source/curescan"
	_ "github.com/EscapeBearSecond/curescan/server/source/example"
	_ "github.com/EscapeBearSecond/curescan/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
