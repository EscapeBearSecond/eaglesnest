package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type StatisticsApi struct {
}

func (s *StatisticsApi) GetVulnsInfo(c *gin.Context) {
	searchResult := &request.SearchJobResult{}
	searchResult.Severity = "critical"
	searchResult.Kind = common.VulnerabilityScan
	searchResult.PageSize = math.MaxInt64
	searchResult.Page = 1
	_, criticalTotal, err := resultService.GetJobResultList(searchResult)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchResult.Severity = "high"
	_, highTotal, err := resultService.GetJobResultList(searchResult)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchResult.Severity = "medium"
	_, mediumTotal, err := resultService.GetJobResultList(searchResult)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchResult.Severity = "low"
	_, lowTotal, err := resultService.GetJobResultList(searchResult)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	var distinctTypeCount int64
	// 查询 kind 为 "2" 的记录，并统计不同的 type 的数量
	err = global.GVA_DB.Model(&curescan.JobResultItem{}).
		Where("kind = ?", common.VulnerabilityScan).
		Select("COUNT(DISTINCT (type))").
		Scan(&distinctTypeCount).Error
	response.OkWithData(gin.H{
		"critical": criticalTotal,
		"high":     highTotal,
		"medium":   mediumTotal,
		"low":      lowTotal,
		"total":    criticalTotal + highTotal + mediumTotal + lowTotal,
		"kindNum":  distinctTypeCount,
	}, c)
}

func (s *StatisticsApi) GetTaskInfo(c *gin.Context) {
	searchTask := request.SearchTask{}
	searchTask.Status = common.Running
	searchTask.TaskPlan = []int{common.ExecuteImmediately, common.ExecuteLater}
	searchTask.Page = 1
	searchTask.PageSize = math.MaxInt64
	_, runningTotal, err := taskService.GetTaskList(searchTask)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchTask.Status = common.Created
	_, createdTotal, err := taskService.GetTaskList(searchTask)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchTask.Status = common.Stopped
	_, stoppedTotal, err := taskService.GetTaskList(searchTask)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchTask.Status = common.Success
	_, successTotal, err := taskService.GetTaskList(searchTask)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	searchTask.Status = common.Failed
	_, failedTotal, err := taskService.GetTaskList(searchTask)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	var distinctIPCount int64

	err = global.GVA_DB.Raw(`SELECT COUNT(DISTINCT (ip))
    FROM (
        SELECT UNNEST(target_ip) AS ip
        FROM cs_task
    ) AS subquery`).Scan(&distinctIPCount).Error
	// err = global.GVA_DB.Model(&curescan.Task{}).
	// 	Select("COUNT(DISTINCT UNNEST(target_ip))").
	// 	Scan(&distinctIPCount).Error
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(gin.H{
		"running":   runningTotal,
		"wait":      createdTotal,
		"stopped":   stoppedTotal,
		"success":   successTotal,
		"failed":    failedTotal,
		"total":     runningTotal + createdTotal + stoppedTotal + successTotal + failedTotal,
		"targetNum": distinctIPCount,
	}, c)
}

func (s *StatisticsApi) CommonVulnTopN(c *gin.Context) {
	nStr := c.DefaultQuery("n", "10")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, err := resultService.CommonVulnTopN(n)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}

func (s *StatisticsApi) AssetTopN(c *gin.Context) {
	nStr := c.DefaultQuery("n", "10")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	list, err := resultService.AssetTopN(n)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
