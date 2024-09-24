package curescan

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/common/request"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/curescan"
	"gorm.io/gorm"
)

type PortScanService struct {
}

func (o *PortScanService) BatchAdd(data []*curescan.PortScan) error {
	return global.GVA_DB.Model(&curescan.PortScan{}).CreateInBatches(data, 100).Error
}

func (o *PortScanService) BatchAddWithTransaction(tx *gorm.DB, data []*curescan.PortScan) error {
	return tx.Model(&curescan.PortScan{}).CreateInBatches(data, 100).Error
}

// ParseFileTo 从传入的os.File类型的文件中解析CSV数据，并将结果转换为curescan.PortScan类型的切片返回
func (o *PortScanService) ParseFileTo(file *os.File) ([]*curescan.PortScan, error) {
	var data []*curescan.PortScan
	reader := csv.NewReader(file)
	_, err := reader.Read()
	//filedCount := len(record)
	var currentIP string
	var ports []int64
	if err != nil {
		return data, err
	}
	records, err := reader.ReadAll()
	for _, record := range records {
		ip := record[0]
		port, _ := strconv.ParseInt(record[1], 10, 64)
		if currentIP == "" {
			currentIP = ip
		} else if currentIP != ip {
			data = append(data, &curescan.PortScan{IP: currentIP, Ports: ports})
			currentIP = ip
			ports = nil
		}
		ports = append(ports, port)
	}
	if currentIP != "" {
		data = append(data, &curescan.PortScan{IP: currentIP, Ports: ports})
	}
	return data, nil
}

// GetInfoList 根据任务ID、分页信息、排序方式和排序方向获取端口扫描信息列表。
// 该方法旨在为用户提供特定任务下的检查结果分页查询能力，并支持结果的升序或降序排列。
// 主要是为了查看某次任务的检查结果功能服务。
func (o *PortScanService) GetInfoList(taskId uint, page request.PageInfo, order string, desc bool) (list []curescan.PortScan, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.PortScan{})
	var data []curescan.PortScan
	if taskId != 0 {
		db = db.Where("task_id = ?", taskId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return data, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool)
		orderMap["ip"] = true
		orderMap["system"] = true
		orderMap["ttl"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return data, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&data).Error
	return data, total, err
}
