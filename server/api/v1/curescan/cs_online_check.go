package curescan

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/response"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/curescan/request"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OnlineCheckApi struct {
}

func (o *OnlineCheckApi) GetInfoList(c *gin.Context) {
	var searchInfo request.SearchInfo
	err := c.ShouldBindJSON(&searchInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(searchInfo.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := onlineCheckService.GetInfoList(searchInfo.TaskID, searchInfo.PageInfo, searchInfo.OrderKey, searchInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchInfo.Page,
		PageSize: searchInfo.PageSize,
	}, "获取成功", c)
}
