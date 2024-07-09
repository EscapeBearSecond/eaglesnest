package initialize

import (
	_ "47.103.136.241/goprojects/curesan/server/source/example"
	_ "47.103.136.241/goprojects/curesan/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
