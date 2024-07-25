package curescan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"47.103.136.241/goprojects/curesan/server/global"
	request2 "47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
	"47.103.136.241/goprojects/curesan/server/model/curescan/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan/response"
	"47.103.136.241/goprojects/eagleeye/pkg/types"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type TaskService struct {
}

var (
	policyService      = &PolicyService{}
	templateService    = &TemplateService{}
	portScanService    = &PortScanService{}
	onlineCheckService = &OnlineCheckService{}
	jobResultService   = &JobResultService{}
)

func (s *TaskService) CreateTask(createTask *request.CreateTask) error {
	// bytes, err := json.Marshal(&createTask.PlanConfig)
	// if err != nil {
	// 	return err
	// }
	var task = curescan.Task{
		TaskName:   createTask.TaskName,
		TaskDesc:   createTask.TaskDesc,
		TaskPlan:   createTask.TaskPlan,
		PlanConfig: createTask.PlanConfig,
		PolicyID:   createTask.PolicyID,
		Status:     createTask.Status,
		TargetIP:   createTask.TargetIP,
	}
	if !errors.Is(global.GVA_DB.Select("task_name").First(&curescan.Task{}, "task_name=?", task.TaskName).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同任务名称，不允许创建")
	}

	err := global.GVA_DB.Create(&task).Error
	if err != nil {
		return err
	}
	// 立即执行
	if createTask.TaskPlan == 1 {
		return s.ExecuteTask(int(task.ID))
	}
	// 稍后执行
	if createTask.TaskPlan == 2 {
		return nil
	}
	// 定时计划
	if createTask.TaskPlan == 3 {
		_, err = global.GVA_Timer.AddTaskByFunc(task.TaskName+time.Now().GoString(), createTask.PlanConfig, func() { s.ExecuteTask(int(task.ID)) }, task.TaskName, cron.WithSeconds())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TaskService) generateCorn(date, timeStr string, frequency int) (string, error) {
	layout := "2006-01-02 15:04:05"
	dateTimeStr := fmt.Sprintf("%s %s", date, timeStr)
	dateTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		return "", fmt.Errorf("解析日期时间出错: %v", err.Error())
	}
	// 生成corn表达式
	var cronExpr string
	switch frequency {
	case 1:
		cronExpr = fmt.Sprintf("0 %d %d %d * ?", dateTime.Second(), dateTime.Minute(), dateTime.Hour(), dateTime.Day())
	case 2:
		dayOfWeek := int(dateTime.Weekday())
		cronExpr = fmt.Sprintf("0 %d %d ? * %d", dateTime.Second(), dateTime.Minute(), dateTime.Hour(), dayOfWeek)
	case 3:
		cronExpr = fmt.Sprintf("0 %d %d * * ?", dateTime.Second(), dateTime.Minute(), dateTime.Hour())
	default:
		return "", fmt.Errorf("不支持的频率: %d", frequency)
	}
	return cronExpr, nil
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
	var onlineConfig request.OnlineConfig
	var portScanConfig request.PortScanConfig
	var jobConfig []request.JobConfig
	if policy.OnlineCheck {
		err = json.Unmarshal([]byte(policy.OnlineConfig), &onlineConfig)
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
		Use:         onlineConfig.Use,
		Timeout:     onlineConfig.Timeout,
		Count:       onlineConfig.Count,
		Format:      onlineConfig.Format,
		Concurrency: onlineConfig.Concurrency,
		RateLimit:   onlineConfig.RateLimit,
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
	jobs, err := s.GenerateJob(id, jobConfig, &taskResult)
	if err != nil {
		return err

	}
	options.Jobs = jobs
	entry, err := global.EagleeyeEngine.NewEntry(options)
	if err != nil {
		return err
	}
	go func() {
		err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			task.Status = 1
			err := s.UpdateTask(task)
			if err != nil {
				return err
			}
			global.GVA_LOG.Info("任务开始执行...", zap.String("taskName", task.TaskName))
			err = entry.Run(context.Background())
			if err != nil {
				return err
			}
			// result := entry.Result()
			// 任务执行成功
			err = portScanService.BatchAdd(taskResult.PortScanList)
			if err != nil {
				return err
			}
			err = onlineCheckService.BatchAdd(taskResult.OnlineCheckList)
			if err != nil {
				return err
			}
			err = jobResultService.BatchAdd(taskResult.JobResultList)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			global.GVA_LOG.Error("任务执行失败", zap.String("taskName", task.TaskName), zap.Error(err))
			task.Status = 3
		} else {
			task.Status = 2
		}
		s.UpdateTask(task)
		global.GVA_LOG.Info("任务执行成功", zap.Int("id", id), zap.String("taskName", task.TaskName))
	}()
	return nil
}

func (s *TaskService) GenerateJob(id int, jobConfig []request.JobConfig, taskResult *response.TaskResult) ([]types.JobOptions, error) {
	jobs := make([]types.JobOptions, len(jobConfig))
	for i, job := range jobConfig {
		// dir := path.Join(global.GVA_CONFIG.AutoCode.Root, "templates", job.Name)
		jobs[i].Name = job.Name
		jobs[i].Kind = job.Kind
		jobs[i].Concurrency = job.Concurrency
		jobs[i].Count = job.Count
		jobs[i].Format = job.Format
		jobs[i].RateLimit = job.RateLimit
		jobs[i].ResultCallback = func(c context.Context, result *types.JobResult) error {
			var data []*curescan.JobResultItem
			for _, item := range result.Items {
				var item = &curescan.JobResultItem{
					Name:             result.Name,
					Kind:             result.Kind,
					TemplateID:       item.EntryID,
					TemplateName:     item.TemplateName,
					Host:             item.Host,
					Type:             item.Type,
					Severity:         item.Severity,
					Port:             item.Port,
					Scheme:           item.Scheme,
					URL:              item.URL,
					Path:             item.Path,
					Matched:          item.Matched,
					ExtractedResults: item.ExtractedResults,
					Description:      item.Description,
					TaskID:           uint(id),
					EntryID:          result.EntryID,
				}
				data = append(data, item)
			}
			taskResult.JobResultList = data
			return nil
		}
		jobs[i].GetTemplates = func() []*types.RawTemplate {
			var rawTemplates []*types.RawTemplate
			templates, err := templateService.GetTemplatesByIds(job.Templates)
			if err != nil {
				return nil
			}
			for _, template := range templates {
				rawTemplates = append(rawTemplates, &types.RawTemplate{
					ID:       template.TemplateName,
					Original: template.TemplateContent,
				})
			}
			return rawTemplates
		}
		return jobs, nil
	}
	return jobs, nil
}
