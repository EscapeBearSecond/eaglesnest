package initialize

import (
	"path"

	"github.com/EscapeBearSecond/eaglesnest/server/global"
	falcon "github.com/EscapeBearSecond/falcon/pkg/sdk"
)

func FalconEngine() {
	var err error
	global.FalconEngine, err = falcon.NewEngine(falcon.WithDirectory(path.Join(global.GVA_CONFIG.AutoCode.Root, "server", "results")))
	if err != nil {
		panic(err)
	}
}
