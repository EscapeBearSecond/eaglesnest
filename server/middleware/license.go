package middleware

import (
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/eagleeye/pkg/license"
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
