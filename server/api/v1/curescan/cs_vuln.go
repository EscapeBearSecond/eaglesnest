package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"47.103.136.241/goprojects/curescan/server/utils"
	"github.com/gin-gonic/gin"
)

type VulnApi struct {
}

func (a *VulnApi) GetVulnList(c *gin.Context) {
	var searchVuln request.SearchVuln
	err := c.ShouldBindJSON(&searchVuln)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(searchVuln.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := vulnService.GetVulnList(&searchVuln)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchVuln.Page,
		PageSize: searchVuln.PageSize,
	}, "获取成功", c)
}
func (a *VulnApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&curescan.Vuln{})
	if err != nil {
		response.FailWithMessage("migrate table failed", c)
		return
	}
	response.OkWithMessage("migrate table success", c)
}
