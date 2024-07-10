package curescan

import (
	"47.103.136.241/goprojects/curesan/server/model/common/response"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"47.103.136.241/goprojects/curesan/server/utils"
	"github.com/gin-gonic/gin"
)

type AreaApi struct {
}

func (csa *AreaApi) CreateArea(c *gin.Context) {
	var createArea request.CreateArea
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
	err = areaService.CreateArea(modelArea)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
