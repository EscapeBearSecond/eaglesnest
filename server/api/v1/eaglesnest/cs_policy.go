package eaglesnest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/common"
	"strconv"

	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/response"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/request"
	csresponse "github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/response"
	"github.com/EscapeBearSecond/eaglesnest/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PolicyApi struct {
}

// MigrateTable
// @deprecated
// 初始化Policy数据库表
func (p *PolicyApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&eaglesnest.Policy{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var po = eaglesnest.Policy{
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
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("创建策略失败!", zap.Any("err", err.Error()))
		}
	}()
	err = c.ShouldBindJSON(&createPolicy)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(createPolicy, utils.CreatePolicyVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !createPolicy.OnlineConfig.Use || !createPolicy.PortScanConfig.Use {
		response.FailWithMessage("请开启在线检测和端口扫描", c)
		return
	}
	if len(createPolicy.PolicyConfig) == 0 && !createPolicy.OnlineConfig.Use && !createPolicy.PortScanConfig.Use {
		err = global.NoDataSelected
		response.FailWithMessage("请至少选择一个扫描类型", c)
		return
	}
	if createPolicy.OnlineConfig.Format == "" {
		createPolicy.OnlineConfig.Format = "csv"
	}
	if createPolicy.PortScanConfig.Format == "" {
		createPolicy.PortScanConfig.Format = "csv"
	}

	for _, policyConfig := range createPolicy.PolicyConfig {
		policyConfig.Name = common.JobTypeName[policyConfig.Kind]
		if policyConfig.IsAll {
			err = global.GVA_DB.Model(&eaglesnest.Template{}).Select("id").Where("template_type", policyConfig.Kind).Scan(&policyConfig.Templates).Error
			if err != nil {
				response.FailWithMessage("全选模板失败", c)
				return
			}
		}
		if policyConfig.Format == "" {
			policyConfig.Format = "excel"
		}

	}
	policyConfig, err := json.Marshal(createPolicy.PolicyConfig)
	if err != nil {
		response.FailWithMessage("任务参数错误", c)
		return
	}

	onlineConfig, err := json.Marshal(createPolicy.OnlineConfig)
	if err != nil {
		response.FailWithMessage("在线检测参数错误", c)
		return
	}
	portScanConfig, err := json.Marshal(createPolicy.PortScanConfig)
	if err != nil {
		response.FailWithMessage("端口扫描参数错误", c)
		return
	}
	scanType := make([]string, len(createPolicy.PolicyConfig))
	templates := make([]int64, 0)
	for i := 0; i < len(createPolicy.PolicyConfig); i++ {
		scanType[i] = createPolicy.PolicyConfig[i].Kind
		templates = append(templates, createPolicy.PolicyConfig[i].Templates...)
	}
	if createPolicy.IgnoredIP == nil {
		createPolicy.IgnoredIP = []string{}
	}
	ip, _ := utils.GetLocalIP()
	createPolicy.IgnoredIP = append(createPolicy.IgnoredIP, ip)
	var modelPolicy = eaglesnest.Policy{
		PolicyName:     createPolicy.PolicyName,
		PolicyDesc:     createPolicy.PolicyDesc,
		ScanType:       scanType,
		PolicyConfig:   string(policyConfig),
		OnlineCheck:    createPolicy.OnlineConfig.Use,
		OnlineConfig:   string(onlineConfig),
		PortScan:       createPolicy.PortScanConfig.Use,
		PortScanConfig: string(portScanConfig),
		Templates:      templates,
		IgnoredIP:      createPolicy.IgnoredIP,
	}
	err = policyService.CreatePolicy(&modelPolicy)
	if err != nil {
		if errors.Is(err, global.HasExisted) {
			response.FailWithMessage("策略已存在", c)
			return
		}
		response.FailWithMessage("创建失败", c)
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
	templates := make([]int64, 0)
	for i := 0; i < len(updatePolicy.PolicyConfig); i++ {
		scanType[i] = updatePolicy.PolicyConfig[i].Kind
		templates = append(templates, updatePolicy.PolicyConfig[i].Templates...)
	}
	var modelPolicy = eaglesnest.Policy{
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
	var onlineConfig csresponse.OnlineConfig
	err = json.Unmarshal([]byte(policy.OnlineConfig), &onlineConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var portScanConfig csresponse.PortScanConfig
	err = json.Unmarshal([]byte(policy.PortScanConfig), &portScanConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var policyConfig []csresponse.JobConfig
	err = json.Unmarshal([]byte(policy.PolicyConfig), &policyConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	policyDetail := &csresponse.PolicyDetail{
		GvaModel:       policy.GvaModel,
		PolicyName:     policy.PolicyName,
		PolicyDesc:     policy.PolicyDesc,
		ScanType:       policy.ScanType,
		PolicyConfig:   policyConfig,
		OnlineCheck:    policy.OnlineCheck,
		OnlineConfig:   onlineConfig,
		PortScan:       policy.PortScan,
		PortScanConfig: portScanConfig,
		Templates:      policy.Templates,
		IgnoredIP:      policy.IgnoredIP,
	}
	response.OkWithData(policyDetail, c)
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
