package initialize

import (
	"path"

	"github.com/EscapeBearSecond/curescan/server/global"
	eagleeye "github.com/EscapeBearSecond/eagleeye/pkg/sdk"
)

func EagleeyeEngine() {
	var err error
	global.EagleeyeEngine, err = eagleeye.NewEngine(eagleeye.WithDirectory(path.Join(global.GVA_CONFIG.AutoCode.Root, "server", "results")))
	if err != nil {
		panic(err)
	}
}
