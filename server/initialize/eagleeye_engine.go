package initialize

import (
	"path"

	"47.103.136.241/goprojects/curesan/server/global"
	eagleeye "47.103.136.241/goprojects/eagleeye/pkg/sdk"
)

func EagleeyeEngine() {
	var err error
	global.EagleeyeEngine, err = eagleeye.NewEngine(eagleeye.WithDirectory(path.Join(global.GVA_CONFIG.AutoCode.Root, "results")))
	if err != nil {
		panic(err)
	}
}
