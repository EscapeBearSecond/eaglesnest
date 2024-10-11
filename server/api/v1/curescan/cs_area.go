package curescan

import (
	"errors"
	"github.com/EscapeBearSecond/curescan/server/global"
	"github.com/EscapeBearSecond/curescan/server/model/common/response"
	"github.com/EscapeBearSecond/curescan/server/model/curescan"
	"github.com/EscapeBearSecond/curescan/server/model/curescan/request"
	"github.com/EscapeBearSecond/curescan/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type AreaApi struct {
}

func (csa *AreaApi) CreateArea(c *gin.Context) {
	var createArea request.CreateArea
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("创建区域失败", zap.Error(err))
		}
	}()
	err = c.ShouldBindJSON(&createArea)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(createArea, utils.CreateAreaVerify)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	var ips = createArea.AreaIP
	err = utils.ValidateIP(ips)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelArea = curescan.Area{
		AreaName: createArea.AreaName,
		AreaIP:   createArea.AreaIP,
		AreaDesc: createArea.AreaDesc,
	}
	modelArea.CreatedBy = utils.GetUserID(c)
	modelArea.UpdatedBy = utils.GetUserID(c)
	err = areaService.CreateArea(&modelArea)
	if err != nil {
		if errors.Is(err, global.HasExisted) {
			response.FailWithMessage("区域已存在", c)
			return
		}
		response.FailWithMessage("创建失败", c)
		return
	}
	response.Ok(c)
}

func (csa *AreaApi) DeleteAreaByID(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("删除区域失败", zap.Error(err))
		}
	}()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = areaService.DeleteArea(id)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.Ok(c)
}

func (csa *AreaApi) UpdateArea(c *gin.Context) {
	var updateArea request.UpdateArea
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("更新区域失败", zap.Error(err))
		}
	}()
	err = c.ShouldBindJSON(&updateArea)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(updateArea, utils.CreateAreaVerify)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	var ips = updateArea.AreaIP
	err = utils.ValidateIP(ips)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelArea = curescan.Area{
		GvaModel: global.GvaModel{
			ID: updateArea.ID,
		},
		AreaName: updateArea.AreaName,
		AreaIP:   updateArea.AreaIP,
		AreaDesc: updateArea.AreaDesc,
	}
	modelArea.UpdatedBy = utils.GetUserID(c)
	err = areaService.UpdateArea(&modelArea)
	if err != nil {
		if errors.Is(err, global.HasExisted) {
			response.FailWithMessage("区域已存在", c)
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (csa *AreaApi) GetAreaById(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("获取区域详情失败", zap.Error(err))
		}
	}()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	area, err := areaService.GetAreaById(id)
	if err != nil {
		if errors.Is(err, global.NoDataFound) {
			response.FailWithMessage("未找到数据", c)
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(area, c)
}

func (csa *AreaApi) GetAreaList(c *gin.Context) {
	var searchArea request.SearchArea
	var err error
	defer func() {
		if err != nil {
			global.GVA_LOG.Error("获取区域列表失败", zap.Error(err))
		}
	}()
	err = c.ShouldBindJSON(&searchArea)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(searchArea.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, total, err := areaService.GetAreaList(&searchArea.Area, searchArea.PageInfo, searchArea.OrderKey, searchArea.Desc)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchArea.Page,
		PageSize: searchArea.PageSize,
	}, "获取成功", c)
}
