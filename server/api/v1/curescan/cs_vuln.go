package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"fmt"
	"github.com/gin-gonic/gin"
)

type VulnApi struct {
}

func (a *VulnApi) GetVulnList(c *gin.Context) {
	vulnList := make([]*curescan.Vuln, 10)
	for i := range 10 {
		vulnList[i] = &curescan.Vuln{
			TemplateID:  fmt.Sprintf("template_id_%d", i),
			Description: fmt.Sprintf("vuln_desc_%d", i),
			Name:        fmt.Sprintf("vuln_name_%d", i),
			Severity:    "high",
			Author:      "author",
			Reference:   "reference",
			Classification: map[string]interface{}{
				"cve":    "cve_id",
				"cwe":    "cwe_id",
				"cvss":   "cvss_id",
				"cvssv3": "cvssv3_id",
			},
			Remediation: "remediation",
		}
	}
	resMap := map[string]interface{}{
		"list":     vulnList,
		"total":    10,
		"page":     1,
		"pageSize": 10,
	}
	response.OkWithData(resMap, c)
}

func (a *VulnApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&curescan.Vuln{})
	if err != nil {
		response.FailWithMessage("migrate table failed", c)
		return
	}
	response.OkWithMessage("migrate table success", c)
}
