package middleware

import (
	"strconv"
	"strings"

	"github.com/EscapeBearSecond/curescan/server/global"
	"github.com/EscapeBearSecond/curescan/server/model/common/response"
	"github.com/EscapeBearSecond/curescan/server/service"
	"github.com/EscapeBearSecond/curescan/server/utils"
	"github.com/gin-gonic/gin"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		waitUse, _ := utils.GetClaims(c)
		// 获取请求的PATH
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, global.GVA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := strconv.Itoa(int(waitUse.AuthorityId))
		c.Set("authorityId", waitUse.AuthorityId)
		e := casbinService.Casbin() // 判断策略中是否存在
		// if e == nil {
		// 	response.FailWithDetailed(gin.H{}, "策略不存在", c)
		// 	c.Abort()
		// 	return
		// }
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
