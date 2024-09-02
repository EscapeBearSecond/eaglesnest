package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	response2 "47.103.136.241/goprojects/curescan/server/model/curescan/response"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"47.103.136.241/goprojects/curescan/server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type StatisticsApi struct {
}

func (s *StatisticsApi) GetVulnsInfo(c *gin.Context) {
	var result response2.SeverityVuln
	var query string
	var err error
	if system.HasAllDataAuthority(c) {
		query = `
        SELECT
            COUNT(*) AS total,
            SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) AS critical,
            SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) AS high,
            SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) AS medium,
            SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) AS low
        FROM (
            SELECT DISTINCT
                host,
                template_id,
                port,
                severity
            FROM cs_job_result
        ) AS distinct_entries
    `
		err = global.GVA_DB.Raw(query).Scan(&result).Error
	} else {
		query = `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) AS critical,
			SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) AS high,
			SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) AS medium,
			SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) AS low
		FROM (
			SELECT DISTINCT
				host,
				template_id,
				port,
				severity
			FROM cs_job_result
			WHERE created_by = ?
		) AS distinct_entries
	`
		err = global.GVA_DB.Raw(query, utils.GetUserID(c)).Scan(&result).Error
	}
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	var distinctTypeCount int64
	// 查询 kind 为 "2" 的记录，并统计不同的 type 的数量
	err = global.GVA_DB.Model(&curescan.JobResultItem{}).
		Where("kind = ?", common.VulnerabilityScan).
		Select("COUNT(DISTINCT (host))").
		Scan(&distinctTypeCount).Error
	response.OkWithData(gin.H{
		"critical": result.Critical,
		"high":     result.High,
		"medium":   result.Medium,
		"low":      result.Low,
		"total":    result.Critical + result.High + result.Medium + result.Low,
		"kindNum":  distinctTypeCount,
	}, c)
}

func (s *StatisticsApi) GetTaskInfo(c *gin.Context) {
	searchTask := request.SearchTask{}
	searchTask.Status = common.Running
	searchTask.TaskPlan = []int{common.ExecuteImmediately, common.ExecuteLater}
	searchTask.Page = 1
	searchTask.PageSize = math.MaxInt64
	searchTask.AllData = system.HasAllDataAuthority(c)
	searchTask.CreatedBy = utils.GetUserID(c)
	fmt.Println("allData", searchTask.AllData)
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
		WHERE created_by = ?
    ) AS subquery`, searchTask.CreatedBy).Scan(&distinctIPCount).Error
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
	list, err := resultService.CommonVulnTopN(n, system.HasAllDataAuthority(c), utils.GetUserID(c))
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
	list, err := resultService.AssetTopN(n, system.HasAllDataAuthority(c), utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithData(list, c)
}
