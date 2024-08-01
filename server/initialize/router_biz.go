package initialize

import (
	"47.103.136.241/goprojects/curescan/server/router"
	"github.com/gin-gonic/gin"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	publicGroup := routers[0]
	privateGroup := routers[1]

	holder(publicGroup, privateGroup)
}
