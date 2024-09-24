package router

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/middleware"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/plugin/email/api"
	"github.com/gin-gonic/gin"
)

type EmailRouter struct{}

func (s *EmailRouter) InitEmailRouter(Router *gin.RouterGroup) {
	emailRouter := Router.Use(middleware.OperationRecord())
	EmailApi := api.ApiGroupApp.EmailApi.EmailTest
	SendEmail := api.ApiGroupApp.EmailApi.SendEmail
	{
		emailRouter.POST("emailTest", EmailApi)  // 发送测试邮件
		emailRouter.POST("sendEmail", SendEmail) // 发送邮件
	}
}
