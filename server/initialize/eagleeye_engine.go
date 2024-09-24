package initialize

import (
	"path"

	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	eagleeye "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/eagleeye/pkg/sdk"
)

func EagleeyeEngine() {
	var err error
	global.EagleeyeEngine, err = eagleeye.NewEngine(eagleeye.WithDirectory(path.Join(global.GVA_CONFIG.AutoCode.Root, "server", "results")))
	if err != nil {
		panic(err)
	}
}
