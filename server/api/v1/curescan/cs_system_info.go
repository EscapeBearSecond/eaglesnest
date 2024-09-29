package curescan

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

type SystemInfoApi struct {
}

func (s *SystemInfoApi) GetSystemInfo(c *gin.Context) {
	info, err := systemInfoService.GetSystemInfo()
	if err != nil {
		response.FailWithMessage("获取系统信息失败", c)
		return
	}
	response.OkWithData(info, c)
}

func (s *SystemInfoApi) UpdateLicense(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("更新证书失败", zap.Error(err))
		}
	}()
	fh, err := c.FormFile("license")
	if err != nil {
		response.FailWithMessage("更新证书失败", c)
		return
	}

	// 检查文件名是否为 license.json
	if fh.Filename != "license.json" {
		response.FailWithMessage("文件名必须是 license.json", c)
		return
	}
	// 定义保存文件的路径
	savePath := "./license.json"

	// 检查文件是否已经存在
	if _, err := os.Stat(savePath); err == nil {
		// 文件存在，删除它
		if err := os.Remove(savePath); err != nil {
			response.FailWithMessage("更新证书失败", c)
			return
		}
	}

	// 保存新上传的文件
	if err := c.SaveUploadedFile(fh, savePath); err != nil {
		response.FailWithMessage("更新证书失败", c)
		return
	}

	response.OkWithMessage("证书更新成功", c)
}
