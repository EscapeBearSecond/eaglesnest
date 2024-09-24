package middleware

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/response"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/eagleeye/pkg/license"
	"github.com/gin-gonic/gin"
)

func LicenseVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		watcher, err := license.Watch("./license.json")
		if err != nil {
			response.FailWithMessage("证书过期", c)
			c.Abort()
			return
		}
		defer watcher.Stop()
		if err = license.Verify(); err != nil {
			response.FailWithMessage("证书过期", c)
			c.Abort()
			return
		}
	}
}
