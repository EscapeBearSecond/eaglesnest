package curescan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"47.103.136.241/goprojects/curesan/server/global"
	request2 "47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan/response"
	"47.103.136.241/goprojects/eagleeye/pkg/types"
	"gorm.io/gorm"
)

type TaskService struct {
}

var (
	policyService      = &PolicyService{}
	templateService    = &TemplateService{}
	portScanService    = &PortScanService{}
	onlineCheckService = &OnlineCheckService{}
)

func (s *TaskService) CreateTask(task *curescan.Task) error {
	if !errors.Is(global.GVA_DB.Select("task_name").First(&curescan.Task{}, "task_name=?", task.TaskName).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同任务名称，不允许创建")
	}

	return global.GVA_DB.Create(&task).Error
}

func (s *TaskService) UpdateTask(task *curescan.Task) error {
	var existingRecord curescan.Task
	err := global.GVA_DB.Select("id", "task_name").Where("task_name=?", task.TaskName).First(&existingRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.GVA_DB.Save(&task).Error
		}
		return err
	}
	if existingRecord.ID != task.ID {
		return errors.New("任务名称已被占用，不允许修改")
	}
	return global.GVA_DB.Save(&task).Error
}

func (s *TaskService) DeleteTask(id int) error {
	return global.GVA_DB.Delete(&curescan.Task{}, id).Error
}

func (s *TaskService) GetTaskById(id int) (*curescan.Task, error) {
	var task curescan.Task
	err := global.GVA_DB.Select("id", "task_name", "task_desc", "status", "target_ip", "policy_id", "task_plan",
		"plan_config", "created_at", "updated_at", "deleted_at").Where("id=?", id).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) GetTaskList(task curescan.Task, page request2.PageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.Task{}).Select("id", "task_name", "task_desc", "status", "target_ip", "policy_id", "task_plan",
		"plan_config", "created_at", "updated_at", "deleted_at")
	var tasks []curescan.Task
	if task.TaskName != "" {
		db = db.Where("task_name LIKE ?", "%"+task.TaskName+"%")
	}
	if task.TaskPlan != 0 {
		db = db.Where("task_plan=?", task.TaskPlan)
	}
	if task.Status != 0 {
		db = db.Where("status=?", task.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return tasks, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		orderMap := make(map[string]bool)
		orderMap["id"] = true
		orderMap["task_name"] = true
		orderMap["status"] = true
		if !orderMap[order] {
			err = fmt.Errorf("非法的排序字段: %s", order)
			return tasks, total, err
		}
		OrderStr = order
		if desc {
			OrderStr += " desc"
		}
	}
	err = db.Order(OrderStr).Find(&tasks).Error
	return tasks, total, err
}

func (s *TaskService) ExecuteTask(id int) error {
	var taskResult response.TaskResult
	task, err := s.GetTaskById(id)
	if err != nil {
		return err
	}
	policy, err := policyService.GetPolicyById(int(task.PolicyID))
	if err != nil {
		return err
	}
	var onlineConfg request.OnlineConfig
	var portScanConfig request.PortScanConfig
	var jobConfig []request.JobConfig
	if policy.OnlineCheck {
		err = json.Unmarshal([]byte(policy.OnlineConfig), &onlineConfg)
		if err != nil {
			return err
		}
	}
	if policy.PortScan {
		err = json.Unmarshal([]byte(policy.PortScanConfig), &portScanConfig)
		if err != nil {
			return err
		}
	}
	if policy.PolicyConfig != "" {
		err = json.Unmarshal([]byte(policy.PolicyConfig), &jobConfig)
		if err != nil {
			return err
		}
	}

	options := &types.Options{}
	options.Targets = task.TargetIP
	options.ExcludeTargets = policy.IgnoredIP
	options.PortScanning = types.PortScanningOptions{
		Use:         portScanConfig.Use,
		Timeout:     portScanConfig.Timeout,
		Count:       portScanConfig.Count,
		Format:      portScanConfig.Format,
		Concurrency: portScanConfig.Concurrency,
		RateLimit:   portScanConfig.RateLimit,
		Ports:       portScanConfig.Ports,
		ResultCallback: func(c context.Context, result *types.PortResult) error {
			if len(result.Items) == 0 {
				return nil
			}
			var data []*curescan.PortScan
			ipPortMap := make(map[string][]int64)

			for _, item := range result.Items {
				ip := item.IP
				port := int64(item.Port)
				ipPortMap[ip] = append(ipPortMap[ip], port)
			}

			for ip, ports := range ipPortMap {
				data = append(data, &curescan.PortScan{
					IP:      ip,
					Ports:   ports,
					TaskID:  uint(id),
					EntryID: result.Items[0].EntryID,
				})
			}
			taskResult.PortScanList = data
			return nil
		},
	}
	options.HostDiscovery = types.HostDiscoveryOptions{
		Use:         onlineConfg.Use,
		Timeout:     onlineConfg.Timeout,
		Count:       onlineConfg.Count,
		Format:      onlineConfg.Format,
		Concurrency: onlineConfg.Concurrency,
		RateLimit:   onlineConfg.RateLimit,
		ResultCallback: func(c context.Context, result *types.PingResult) error {
			var data []*curescan.OnlineCheck
			for _, result := range result.Items {
				data = append(data, &curescan.OnlineCheck{
					IP:      result.IP,
					Active:  result.Active,
					System:  result.OS,
					TTL:     result.TTL,
					TaskID:  uint(id),
					EntryID: result.EntryID,
				})
			}
			taskResult.OnlineCheckList = data
			return nil
		},
	}
	jobs := make([]types.JobOptions, len(jobConfig))
	for i, job := range jobConfig {
		jobs[i].Name = job.Name
		jobs[i].Kind = job.Kind
		jobs[i].Concurrency = job.Concurrency
		jobs[i].Count = job.Count
		jobs[i].Format = job.Format
		jobs[i].RateLimit = job.RateLimit
		jobs[i].ResultCallback = func(c context.Context, result *types.JobResult) error {
			params := []string{}
			for _, result := range result.Items {
				params = append(params, fmt.Sprintf("%#v\n", result))
			}
			fmt.Printf("漏洞扫描: %v\n", params)
			return nil
		}
		// TODO: 考虑并发
		for _, template := range job.Templates {
			template, err := templateService.GetTemplateById(int(template))
			if err != nil {
				continue
			}
			// 创建文件夹和文件
			dir := path.Join(global.GVA_CONFIG.AutoCode.Root, "templates", job.Name)
			_, err = os.Stat(dir)
			if os.IsNotExist(err) {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					return fmt.Errorf("加载模板出错: %s", err.Error())
				}
			}
			filePath := path.Join(dir, template.TemplateName+".yaml")
			_, err = os.Stat(filePath)
			if os.IsNotExist(err) {
				file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
				if err != nil {
					return fmt.Errorf("加载模板出错: %s", err.Error())
				}
				defer file.Close()
				_, err = file.Write([]byte(template.TemplateContent))
				if err != nil {
					return fmt.Errorf("加载模板出错: %s", err.Error())
				}
			}
		}
	}
	options.Jobs = jobs
	entry, err := global.EagleeyeEngine.NewEntry(options)
	if err != nil {
		return err
	}
	go func() {
		task.Status = 1
		err := s.UpdateTask(task)
		if err != nil {
			fmt.Println("更新任务状态时出错: ", err.Error())
		}
		err = entry.Run(context.Background())
		if err != nil {
			fmt.Printf("任务执行出错: %s\n", err.Error())

			task.Status = 3

		}
		// 任务执行成功
		err = portScanService.BatchAdd(taskResult.PortScanList)
		if err != nil {
			fmt.Println("端口扫描-插入数据库时出错: ", err.Error())
		}
		err = onlineCheckService.BatchAdd(taskResult.OnlineCheckList)
		if err != nil {
			fmt.Println("在线检测-插入数据库时出错: ", err.Error())
		}
		task.Status = 2
		err = s.UpdateTask(task)
		if err != nil {
			fmt.Println("更新任务状态时出错: ", err.Error())
		}

	}()
	return nil
}
