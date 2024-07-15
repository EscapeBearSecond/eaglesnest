package curescan

import (
	"encoding/json"
	"fmt"
	"strconv"

	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/response"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"47.103.136.241/goprojects/curesan/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PolicyApi struct {
}

// MigrateTable
// @deprecated
// 初始化Policy数据库表
func (p *PolicyApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&curescan.Policy{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var po = curescan.Policy{
		PolicyName:     "",
		PolicyDesc:     "",
		ScanType:       nil,
		PolicyConfig:   "",
		OnlineCheck:    false,
		OnlineConfig:   "",
		PortScan:       false,
		PortScanConfig: "",
		Templates:      nil,
	}
	b, _ := json.Marshal(po)
	fmt.Println(string(b))
	response.Ok(c)
}

func (p *PolicyApi) CreatePolicy(c *gin.Context) {
	var createPolicy request.CreatePolicy
	err := c.ShouldBindJSON(&createPolicy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(createPolicy, utils.CreatePolicyVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelPolicy = curescan.Policy{
		PolicyName:     createPolicy.PolicyName,
		PolicyDesc:     createPolicy.PolicyDesc,
		ScanType:       createPolicy.ScanType,
		PolicyConfig:   createPolicy.PolicyConfig,
		OnlineCheck:    createPolicy.OnlineCheck,
		OnlineConfig:   createPolicy.OnlineConfig,
		PortScan:       createPolicy.PortScan,
		PortScanConfig: createPolicy.PortScanConfig,
		Templates:      createPolicy.Templates,
		IgnoredIP:      createPolicy.IgnoredIP,
	}
	err = policyService.CreatePolicy(&modelPolicy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (p *PolicyApi) UpdatePolicy(c *gin.Context) {
	var updatePolicy request.UpdatePolicy
	err := c.ShouldBindJSON(&updatePolicy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(updatePolicy.CreatePolicy, utils.CreatePolicyVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelPolicy = curescan.Policy{
		GvaModel: global.GvaModel{
			ID: updatePolicy.ID,
		},
		PolicyName:     updatePolicy.PolicyName,
		PolicyDesc:     updatePolicy.PolicyDesc,
		ScanType:       updatePolicy.ScanType,
		PolicyConfig:   updatePolicy.PolicyConfig,
		OnlineCheck:    updatePolicy.OnlineCheck,
		OnlineConfig:   updatePolicy.OnlineConfig,
		PortScan:       updatePolicy.PortScan,
		PortScanConfig: updatePolicy.PortScanConfig,
		Templates:      updatePolicy.Templates,
	}
	err = policyService.UpdatePolicy(&modelPolicy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (p *PolicyApi) DeletePolicy(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = policyService.DeletePolicy(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (p *PolicyApi) GetPolicyById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	policy, err := policyService.GetPolicyById(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(policy, c)

}

func (p *PolicyApi) GetPolicyList(c *gin.Context) {
	var searchPolicy request.SearchPolicy
	err := c.ShouldBindJSON(&searchPolicy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(searchPolicy.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := policyService.GetPolicyList(searchPolicy.Policy, searchPolicy.PageInfo, searchPolicy.OrderKey, searchPolicy.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchPolicy.Page,
		PageSize: searchPolicy.PageSize,
	}, "获取成功", c)
}
