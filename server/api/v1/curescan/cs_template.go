package curescan

import (
	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/response"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"github.com/gin-gonic/gin"
)

type TemplateApi struct {
}

func (t *TemplateApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&curescan.Template{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)

}
func (t *TemplateApi) CreateTemplate(c *gin.Context) {

}
