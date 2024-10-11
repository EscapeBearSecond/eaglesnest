package middleware

import (
	"github.com/EscapeBearSecond/curescan/server/model/common/response"
	"github.com/EscapeBearSecond/eagleeye/pkg/license"
	"github.com/gin-gonic/gin"
)

func LicenseVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		watcher, err := license.Watch("./license.json")
		if err != nil {
			response.LicenseExpired("证书过期", c)
			c.Abort()
			return
		}
		defer watcher.Stop()
		if err = license.Verify(); err != nil {
			response.LicenseExpired("证书过期", c)
			c.Abort()
			return
		}
	}
}
