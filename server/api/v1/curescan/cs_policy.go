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
	policyConfig, _ := json.Marshal(createPolicy.PolicyConfig)
	onlineConfg, _ := json.Marshal(createPolicy.OnlineConfig)
	portScanConfg, _ := json.Marshal(createPolicy.PortScanConfig)
	scanType := make([]string, len(createPolicy.PolicyConfig))
	templates := make([]string, 0)
	for i := 0; i < len(createPolicy.PolicyConfig); i++ {
		scanType[i] = createPolicy.PolicyConfig[i].Kind
		templates = append(templates, createPolicy.PolicyConfig[i].Templates...)
	}
	var modelPolicy = curescan.Policy{
		PolicyName:     createPolicy.PolicyName,
		PolicyDesc:     createPolicy.PolicyDesc,
		ScanType:       scanType,
		PolicyConfig:   string(policyConfig),
		OnlineCheck:    createPolicy.OnlineConfig.Use,
		OnlineConfig:   string(onlineConfg),
		PortScan:       createPolicy.PortScanConfig.Use,
		PortScanConfig: string(portScanConfg),
		Templates:      templates,
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
	policyConfig, _ := json.Marshal(updatePolicy.PolicyConfig)
	onlineConfg, _ := json.Marshal(updatePolicy.OnlineConfig)
	portScanConfg, _ := json.Marshal(updatePolicy.PortScanConfig)
	scanType := make([]string, len(updatePolicy.PolicyConfig))
	templates := make([]string, 0)
	for i := 0; i < len(updatePolicy.PolicyConfig); i++ {
		scanType[i] = updatePolicy.PolicyConfig[i].Kind
		templates = append(templates, updatePolicy.PolicyConfig[i].Templates...)
	}
	var modelPolicy = curescan.Policy{
		GvaModel: global.GvaModel{
			ID: updatePolicy.ID,
		},
		PolicyName:     updatePolicy.PolicyName,
		PolicyDesc:     updatePolicy.PolicyDesc,
		ScanType:       scanType,
		PolicyConfig:   string(policyConfig),
		OnlineCheck:    updatePolicy.OnlineConfig.Use,
		OnlineConfig:   string(onlineConfg),
		PortScan:       updatePolicy.PortScanConfig.Use,
		PortScanConfig: string(portScanConfg),
		Templates:      templates,
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
