package eaglesnest

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"gorm.io/gorm"
)

type OnlineCheckService struct {
}

func (o *OnlineCheckService) BatchAdd(data []*eaglesnest.OnlineCheck) error {
	return global.GVA_DB.Model(&eaglesnest.OnlineCheck{}).CreateInBatches(data, 100).Error
}

func (o *OnlineCheckService) BatchAddWithTransaction(tx *gorm.DB, data []*eaglesnest.OnlineCheck) error {
	return tx.Model(&eaglesnest.OnlineCheck{}).CreateInBatches(data, 100).Error
}

// ParseFileTo 从传入的os.File类型的文件中解析CSV数据，并将结果转换为eaglesnest.OnlineCheck类型的切片返回
func (o *OnlineCheckService) ParseFileTo(file *os.File) ([]*eaglesnest.OnlineCheck, error) {
	var data []*eaglesnest.OnlineCheck
	reader := csv.NewReader(file)
	record, err := reader.Read()
	filedCount := len(record)
	if err != nil {
		return data, err
	}
	record, err = reader.Read()
	for err == nil || errors.Is(err, csv.ErrFieldCount) {
		one := &eaglesnest.OnlineCheck{}
		one.IP = record[0]
		if record[1] == "是" {
			one.Active = true
		} else {
			one.Active = false
		}
		if len(record) == filedCount {
			one.System = record[2]
			ttl, _ := strconv.ParseInt(record[3], 10, 64)
			one.TTL = int(ttl)

		} else {
			one.System = ""
			one.TTL = 0
		}
		data = append(data, one)
		record, err = reader.Read()
	}
	return data, nil
}

// GetInfoList 根据任务ID、分页信息、排序方式和排序方向获取线上检查信息列表。
// 该方法旨在为用户提供特定任务下的检查结果分页查询能力，并支持结果的升序或降序排列。
// 主要是为了查看某次任务的检查结果功能服务。
func (o *OnlineCheckService) GetInfoList(taskId uint, page request.PageInfo, order string, desc bool) (list []eaglesnest.OnlineCheck, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&eaglesnest.OnlineCheck{})
	var data []eaglesnest.OnlineCheck
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
