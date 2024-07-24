package curescan

import (
	"strconv"

	"47.103.136.241/goprojects/curesan/server/global"
	"47.103.136.241/goprojects/curesan/server/model/common/response"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"47.103.136.241/goprojects/curesan/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AreaApi struct {
}

func (csa *AreaApi) CreateArea(c *gin.Context) {
	var createArea request.CreateArea
	// err := utils.BindAndValid(c, &createArea)
	// if err != nil {
	// 	response.FailWithMessage(err.Error(), c)
	// 	return
	// }
	err := c.ShouldBindJSON(&createArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(createArea, utils.CreateAreaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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
	err = areaService.CreateArea(&modelArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (csa *AreaApi) DeleteAreaByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = areaService.DeleteArea(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (csa *AreaApi) UpdateArea(c *gin.Context) {
	var updateArea request.UpdateArea
	err := c.ShouldBindJSON(&updateArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(updateArea, utils.CreateAreaVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
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
	err = areaService.UpdateArea(&modelArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (csa *AreaApi) GetAreaById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	area, err := areaService.GetAreaById(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(area, c)
}

func (csa *AreaApi) GetAreaList(c *gin.Context) {
	var searchArea request.SearchArea
	err := c.ShouldBindJSON(&searchArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(searchArea.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := areaService.GetAreaList(searchArea.Area, searchArea.PageInfo, searchArea.OrderKey, searchArea.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
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
