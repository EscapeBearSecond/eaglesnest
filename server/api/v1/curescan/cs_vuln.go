package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"47.103.136.241/goprojects/curescan/server/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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

type JSONB map[string]interface{}
type ddd struct {
	TemplateID     string   `json:"template_id" gorm:"column:template_id;type:text;uniqueIndex;comment:模板id" json:"templateId"`
	Name           string   `json:"name" gorm:"column:name;type:text;uniqueIndex;not null;comment:漏洞名称" json:"name"`
	Author         string   `json:"author" gorm:"column:author;type:text;not null;comment:漏洞作者" json:"author"`
	Severity       string   `json:"severity" gorm:"column:severity;type:text;not null;comment:漏洞等级" json:"severity"`
	Description    string   `json:"description" gorm:"column:description;type:text;comment:漏洞描述" json:"description"`
	Reference      []string `json:"reference" gorm:"column:reference;type:json;comment:引用信息" json:"reference"`
	Classification JSONB    `json:"classification" gorm:"column:classification;type:jsonb;comment:其他分类信息" json:"classification"`
	Remediation    string   `json:"remediation" gorm:"column:remediation;type:text;comment:修复建议" json:"remediation"`
}

func (a *VulnApi) DataMod(c *gin.Context) {
	// 读取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	// 解析 JSON 文件内容
	var vulns []*ddd
	if err := json.NewDecoder(file).Decode(&vulns); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
		return
	}
	// 补充缺失字段
	for i := range vulns {
		if vulns[i].TemplateID == "" {
			vulns[i].TemplateID = ""
		}
		if vulns[i].Name == "" {
			vulns[i].Name = ""
		}
		if vulns[i].Author == "" {
			vulns[i].Author = ""
		}
		if vulns[i].Severity == "" {
			vulns[i].Severity = ""
		}
		if vulns[i].Description == "" {
			vulns[i].Description = ""
		}
		if vulns[i].Reference == nil {
			vulns[i].Reference = []string{}
		}
		if vulns[i].Classification == nil {
			vulns[i].Classification = map[string]interface{}{}
		}
		if vulns[i].Remediation == "" {
			vulns[i].Remediation = ""
		}
	}
	// 将补充后的数据写入文件
	fileName := "processed_vulns.json"
	fileData, err := json.MarshalIndent(vulns, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	if err := os.WriteFile(fileName, fileData, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	// 返回文件下载
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/json")
	c.File(fileName)

	// 删除文件（可选，确保不会保留临时文件）
	if err := os.Remove(fileName); err != nil {
		fmt.Println("Failed to remove file:", err)
	}
}
