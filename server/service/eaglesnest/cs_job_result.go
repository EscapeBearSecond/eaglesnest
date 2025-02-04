package eaglesnest

import (
	"fmt"
	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/common"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/response"
	"gorm.io/gorm"
)

type JobResultService struct {
}

func (j *JobResultService) BatchAdd(data []*eaglesnest.JobResultItem) error {
	return global.GVA_DB.Model(&eaglesnest.JobResultItem{}).CreateInBatches(data, 100).Error
}

func (j *JobResultService) BatchAddWithTransaction(tx *gorm.DB, data []*eaglesnest.JobResultItem) error {
	return tx.Transaction(func(tx1 *gorm.DB) error {
		return tx1.Model(&eaglesnest.JobResultItem{}).CreateInBatches(data, 100).Error
	})
}

func (j *JobResultService) GetJobResultList(info *request.SearchJobResult) (list interface{}, total int64, err error) {
	jobResult := info.JobResultItem
	page := info.PageInfo
	order := info.OrderKey
	desc := info.Desc
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	// distinctFields := info.DistinctFields
	db := global.GVA_DB.Model(&eaglesnest.JobResultItem{})
	var jobResultList []*eaglesnest.JobResultItem
	if jobResult.Kind != "" {
		db = db.Where("kind = ?", jobResult.Kind)
	}
	if jobResult.Name != "" {
		db = db.Where("name = ?", jobResult.Name)
	}
	if jobResult.Severity != "" {
		db = db.Where("severity = ?", jobResult.Severity)
	}
	// 去重处理
	// if len(distinctFields) > 0 {
	// 	db = db.Select(distinctFields).Distinct(distinctFields)
	// }

	err = db.Count(&total).Error
	if err != nil {
		return jobResultList, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool)
		orderMap["id"] = true
		orderMap["kind"] = true
		orderMap["name"] = true
		orderMap["severity"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return jobResultList, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&jobResultList).Error
	return jobResultList, total, err
}

func (j *JobResultService) CommonVulnTopN(n int, allData bool, userId uint) (interface{}, error) {
	var res = make([]*response.VulnTopN, 0)
	var err error
	var query string
	if allData {
		query = `
        SELECT
            distinct_entries.template_name AS name,
            distinct_entries.severity,
            COUNT(*) AS count
        FROM (
            SELECT DISTINCT
                template_name,
                severity,
                host,
                port
            FROM cs_job_result
			WHERE kind = ?
        ) AS distinct_entries
        GROUP BY
            distinct_entries.template_name,
            distinct_entries.severity
        ORDER BY
            count DESC
        LIMIT 10
    `
		// Execute the query
		err = global.GVA_DB.Raw(query, common.VulnerabilityScan).Scan(&res).Error
	} else {
		query = `
        SELECT
            distinct_entries.template_name AS name,
            distinct_entries.severity,
            COUNT(*) AS count
        FROM (
            SELECT DISTINCT
                template_name,
                severity,
                host,
                port
            FROM cs_job_result
			WHERE kind = ?
			AND created_by=?
        ) AS distinct_entries
        GROUP BY
            distinct_entries.template_name,
            distinct_entries.severity
        ORDER BY
            count DESC
        LIMIT 10
    `
		// Execute the query
		err = global.GVA_DB.Raw(query, common.VulnerabilityScan, userId).Scan(&res).Error
	}

	// err := global.GVA_DB.Table("cs_job_result").
	// 	Select("template_name as name, severity, COUNT(*) as count").
	// 	Joins("JOIN (SELECT DISTINCT template_name, severity, host, port FROM cs_job_result) AS distinct_entries ON cs_job_result.template_name = distinct_entries.template_name AND cs_job_result.severity = distinct_entries.severity AND cs_job_result.host = distinct_entries.host AND cs_job_result.port = distinct_entries.port").
	// 	Group("template_name, severity").
	// 	Order("count DESC").
	// 	Limit(n).
	// 	Scan(&res).Error
	return res, err
}

func (j *JobResultService) AssetTopN(n int, allData bool, userId uint) (interface{}, error) {
	var results = make([]*response.AssetTopN, 0)
	var query string
	var err error
	if allData {
		query = `
		SELECT
			host,
			COUNT(*) AS count,
			SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) AS critical,
			SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) AS high,
			SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) AS medium,
			SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) AS low
		FROM (
			SELECT 
				host,
				severity
			FROM 
				(SELECT DISTINCT host, template_id, port, severity FROM cs_job_result WHERE kind = ?) AS distinct_entries
		) AS severity_entries
		GROUP BY
			host
		ORDER BY
			count DESC
		LIMIT 10
    `
		// Execute the query
		err = global.GVA_DB.Raw(query, common.VulnerabilityScan).Scan(&results).Error
	} else {
		query = `
		SELECT
			host,
			COUNT(*) AS count,
			SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) AS critical,
			SUM(CASE WHEN severity = 'high' THEN 1 ELSE 0 END) AS high,
			SUM(CASE WHEN severity = 'medium' THEN 1 ELSE 0 END) AS medium,
			SUM(CASE WHEN severity = 'low' THEN 1 ELSE 0 END) AS low
		FROM (
			SELECT 
				host,
				severity
			FROM 
				(SELECT DISTINCT host, template_id, port, severity FROM cs_job_result WHERE kind = ? AND created_by = ?) AS distinct_entries
		) AS severity_entries
		GROUP BY
			host
		ORDER BY
			count DESC
		LIMIT 10
    `
		// Execute the query
		err = global.GVA_DB.Raw(query, common.VulnerabilityScan, userId).Scan(&results).Error
	}

	return results, err
}
