package curescan

import (
	"fmt"
	"strconv"

	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/response"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"47.103.136.241/goprojects/curesan/server/utils"
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

// CreateTemplate
// 创建模板
func (t *TemplateApi) CreateTemplate(c *gin.Context) {
	var createTemplate request.CreateTemplate
	err := c.ShouldBindJSON(&createTemplate)
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(createTemplate, utils.CreateTemplateVerify)
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelTemplate = curescan.Template{
		TemplateName:    createTemplate.TemplateName,
		TemplateType:    createTemplate.TemplateType,
		TemplateDesc:    createTemplate.TemplateDesc,
		TemplateContent: createTemplate.TemplateContent,
	}

	err = templateService.CreateTemplate(&modelTemplate)
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// ImportTemplateContent
// 导入模板内容
func (t *TemplateApi) ImportTemplateContent(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var buffer = make([]byte, fileHeader.Size)
	n, err := file.Read(buffer)
	fmt.Println("文件长度：", n)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if int64(n) != fileHeader.Size {
		response.FailWithMessage(fmt.Sprintf("文件缺失, 源文件大小: %d, 读取大小: %d ", fileHeader.Size, n), c)
		return
	}
	response.OkWithData(string(buffer), c)
}

func (t *TemplateApi) DeleteTemplate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = templateService.DeleteTemplate(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

// GetTemplateById
// 返回模板以及模板内容
func (t *TemplateApi) GetTemplateById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	template, err := templateService.GetTemplateById(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(template, c)
}

// GetTemplateList
// 返回分页排序后的模板列表，不返回模板内容
func (t *TemplateApi) GetTemplateList(c *gin.Context) {
	var searchTemplate request.SearchTemplate
	err := c.ShouldBindJSON(&searchTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(searchTemplate.PageInfo)
	err = utils.Verify(searchTemplate.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := templateService.GetTemplateList(searchTemplate.Template, searchTemplate.PageInfo, searchTemplate.OrderKey, searchTemplate.Desc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchTemplate.Page,
		PageSize: searchTemplate.PageSize,
	}, "获取成功", c)
}

func (t *TemplateApi) UpdateTemplate(c *gin.Context) {
	var updateTemplate request.UpdateTemplate
	err := c.ShouldBindJSON(&updateTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(updateTemplate, utils.CreateTemplateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelTemplate = curescan.Template{
		GvaModel: global.GvaModel{
			ID: updateTemplate.ID,
		},
		TemplateName:    updateTemplate.TemplateName,
		TemplateType:    updateTemplate.TemplateType,
		TemplateDesc:    updateTemplate.TemplateDesc,
		TemplateContent: updateTemplate.TemplateContent,
	}
	err = templateService.UpdateTemplate(&modelTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
