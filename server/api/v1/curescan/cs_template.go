package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"47.103.136.241/goprojects/curescan/server/utils"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type TemplateApi struct {
}

func (t *TemplateApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&curescan.Template{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.AutoMigrate(&curescan.Area{})
	err = global.GVA_DB.AutoMigrate(&curescan.JobResultItem{})
	response.Ok(c)
}

// CreateTemplate
// 创建模板
func (t *TemplateApi) CreateTemplate(c *gin.Context) {
	var createTemplate request.CreateTemplate
	err := c.ShouldBindJSON(&createTemplate)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.String("url", c.Request.URL.Path), zap.String("error", err.Error()))
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(createTemplate, utils.CreateTemplateVerify)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.String("url", c.Request.URL.Path), zap.String("error", err.Error()))
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelTemplate = curescan.Template{
		TemplateType:    createTemplate.TemplateType,
		TemplateContent: createTemplate.TemplateContent,
	}
	err = templateService.ParseTemplateContent(&modelTemplate)
	if err != nil {
		global.GVA_LOG.Error("模板内容解析失败", zap.String("uri", c.Request.RequestURI), zap.String("error", err.Error()))
		response.FailWithMessage(err.Error(), c)
		return
	}
	modelTemplate.CreatedBy = utils.GetUserID(c)
	modelTemplate.UpdatedBy = utils.GetUserID(c)
	err = templateService.CreateTemplate(&modelTemplate)
	if err != nil {
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

// gzipWriter 结构体
type gzipWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
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
	err = utils.Verify(searchTemplate.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 设置响应头，指明内容使用了gzip压缩
	c.Writer.Header().Set("Content-Encoding", "gzip")
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)

	// 创建gzip.Writer
	gz := gzip.NewWriter(c.Writer)
	defer gz.Close()

	// 使用 gzipWriter 进行压缩写入
	gzipWriter := &gzipWriter{Writer: gz, ResponseWriter: c.Writer}

	list, total, err := templateService.GetTemplateList(searchTemplate)
	if err != nil {
		response.FailWithMessage("获取数据失败", c)
		return
	}

	// 构建响应结果
	result := response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchTemplate.Page,
		PageSize: searchTemplate.PageSize,
	}
	mapRest := make(map[string]interface{})
	mapRest["data"] = result
	mapRest["code"] = response.SUCCESS
	mapRest["msg"] = "查询成功"

	// 编码并写入数据到gzipWriter

	if err := json.NewEncoder(gzipWriter).Encode(mapRest); err != nil {
		response.FailWithMessage("编解码错误", c)
		return
	}
	gzipWriter.Flush()
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
		TemplateType:    updateTemplate.TemplateType,
		TemplateContent: updateTemplate.TemplateContent,
	}
	err = templateService.ParseTemplateContent(&modelTemplate)
	if err != nil {
		global.GVA_LOG.Error("模板内容解析失败", zap.String("uri", c.Request.RequestURI), zap.String("error", err.Error()))
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = templateService.UpdateTemplate(&modelTemplate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (t *TemplateApi) ImportTemplates(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		global.GVA_LOG.Error("获取文件失败", zap.String("uri", c.Request.RequestURI), zap.String("error", err.Error()))
		response.FailWithMessage("获取文件失败", c)
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		response.FailWithMessage("没有上传文件", c)
		return
	}
	errorStrings := make([]string, 0)
	types := form.Value["templateType"]
	templateType := types[0]
	templates := make([]*curescan.Template, 0)
	for _, fh := range files {

		file, err := fh.Open()
		if err != nil {
			errorStrings = append(errorStrings, fmt.Sprintf("open file [%s] error, err: [%s]", fh.Filename, err.Error()))
			file.Close()
			continue
		}
		var buffer = make([]byte, fh.Size)
		n, err := file.Read(buffer)
		if err != nil {
			errorStrings = append(errorStrings, fmt.Sprintf("read file [%s] error, err: [%s]", fh.Filename, err.Error()))
			file.Close()
			continue
		}
		if int64(n) != fh.Size {
			errorStrings = append(errorStrings, fmt.Sprintf("file [%s] missing, source file size: %d, read size: %d", fh.Filename, fh.Size, n))
			file.Close()
			continue
		}
		file.Close()

		template := &curescan.Template{}
		template.TemplateContent = string(buffer)
		template.TemplateType = templateType
		err = templateService.ParseTemplateContent(template)
		if err != nil {
			global.GVA_LOG.Error("模板内容解析失败", zap.String("模板id", template.TemplateId))
			errorStrings = append(errorStrings, fmt.Sprintf("parse template [%s] error, err: [%s]", fh.Filename, err.Error()))
			continue
		}
		templates = append(templates, template)
	}
	err = templateService.BatchAdd(templates)
	if err != nil {
		global.GVA_LOG.Error("上传模板失败", zap.String("uri", c.Request.RequestURI), zap.String("error", err.Error()))
		response.FailWithDetailed(errorStrings, "上传模板失败", c)
		return
	}
	if len(errorStrings) > 0 {
		response.OkWithDetailed(errorStrings, "部分模板上传失败", c)
		return
	}
	response.Ok(c)

}

func (t *TemplateApi) TemplateTags(c *gin.Context) {
	var tag1s []string
	var tag2s []string
	var tag3s []string
	var tag4s []string
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		template := curescan.Template{}
		err := tx.Model(template).Select("tag1").Where("tag1 IS NOT NULL AND tag1 != ''").Distinct("tag1").Find(&tag1s).Error
		if err != nil {
			return err
		}
		err = tx.Model(template).Select("tag2").Where("tag2 IS NOT NULL AND tag2 != ''").Distinct("tag2").Find(&tag2s).Error
		if err != nil {
			return err
		}
		err = tx.Model(template).Select("tag3").Where("tag2 IS NOT NULL AND tag3 != ''").Distinct("tag3").Find(&tag3s).Error
		if err != nil {
			return err
		}
		err = tx.Model(template).Select("tag4").Where("tag3 IS NOT NULL AND tag4 != ''").Distinct("tag4").Find(&tag4s).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		global.GVA_LOG.Error("查询数据库出现异常", zap.String("uri", c.Request.RequestURI), zap.String("error", err.Error()))
		response.FailWithMessage("查询数据库出现异常", c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"tag1": tag1s,
		"tag2": tag2s,
		"tag3": tag3s,
		"tag4": tag4s,
	}, c)
}

func (t *TemplateApi) UploadFromZip(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("上传模板失败", zap.String("uri", c.Request.RequestURI), zap.String("error", err.Error()))
		}
		// 删除 zip 文件和解压后的文件夹
		zipExist := utils.FileExist("template.zip")
		if zipExist {
			os.Remove("template.zip")
		}
		dirExist, err := utils.PathExists("template")
		if err == nil {
			if dirExist {
				os.RemoveAll("template")
			}
		}

	}()
	fh, err := c.FormFile("file")
	if err != nil || fh == nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	file, err := fh.Open()
	defer file.Close()
	if err != nil {
		response.FailWithMessage("打开加密文件失败", c)
		return
	}
	encipheredData, err := io.ReadAll(file)
	if err != nil {
		response.FailWithMessage("读取加密文件失败", c)
		return
	}
	err = utils.DecryptFile(encipheredData, "fNwVcQpR", "template.zip")
	if err != nil {
		response.FailWithMessage("解密文件失败", c)
		return
	}
	paths, err := utils.Unzip("template.zip", "template")
	if err != nil {
		response.FailWithMessage("解压文件失败", c)
		return
	}
	templates := make([]*curescan.Template, 0)

	for _, path := range paths {
		if !utils.IsFile(path) {
			continue
		}
		var templateType string
		dir := filepath.Base(filepath.Dir(path))
		if dir == "vuln" || dir == "version_vuln" {
			templateType = common.VulnerabilityScan
		} else if dir == "asset" {
			templateType = common.AssetDiscovery
		} else if dir == "weak" {
			templateType = common.WeakPassword
		} else {
			continue
		}
		file, err := os.Open(path)
		if err != nil {
			response.FailWithMessage("打开文件失败", c)
			return
		}

		buffer, err := os.ReadFile(path)
		if err != nil {
			file.Close()
			response.FailWithMessage("读取文件失败", c)
			return
		}

		template := &curescan.Template{}
		template.TemplateContent = string(buffer)
		template.TemplateType = templateType

		err = templateService.ParseTemplateContent(template)
		if err != nil {
			file.Close()
			response.FailWithMessage(fmt.Sprintf("解析模板%s失败", path), c)

		}
		templates = append(templates, template)
		file.Close()
	}

	err = templateService.BatchAdd(templates)
	if err != nil {
		response.FailWithMessage("添加模板失败", c)
		return
	}
	response.Ok(c)
}

func (t *TemplateApi) LLL(c *gin.Context) {
	templates, err := templateService.GetTemplatesByIds([]int64{130})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(templates, c)
}
