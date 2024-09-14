package curescan

import (
	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"47.103.136.241/goprojects/curescan/server/model/curescan/response"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"47.103.136.241/goprojects/curescan/server/utils"
	"47.103.136.241/goprojects/eagleeye/pkg/report"
	eagleeye "47.103.136.241/goprojects/eagleeye/pkg/sdk"
	"47.103.136.241/goprojects/eagleeye/pkg/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type TaskService struct {
}

var (
	policyService      = &PolicyService{}
	templateService    = &TemplateService{}
	portScanService    = &PortScanService{}
	onlineCheckService = &OnlineCheckService{}
	jobResultService   = &JobResultService{}
	userService        = &system.UserService{}
	assetService       = &AssetService{}
	casbinService      = &system.CasbinService{}
)

var portAssetMap = map[int64]*curescan.Asset{
	3306: {AssetName: "MySQL", AssetType: "服务器设备", AssetModel: "MySQL", SystemType: "MySQL", Manufacturer: "MySQL"},
	3389: {AssetName: "Windows远程桌面", AssetType: "服务器设备", AssetModel: "MySQL", SystemType: "Windows", Manufacturer: "Microsoft"},
	23:   {AssetName: "telnet", AssetType: "网络设备", AssetModel: "telnet", SystemType: "telnet", Manufacturer: "telnet"},
	554:  {AssetName: "rtsp", AssetType: "视频设备", AssetModel: "rtsp", SystemType: "rtsp", Manufacturer: "rtsp"},
}

// 定义端口优先级
var portPriority = map[int64]int64{
	3306: 1,
	3389: 2,
	23:   3,
	554:  4,
}

// 用于存储每个IP的最高优先级端口
var ipPortMap = make(map[string]int64)

var ipPorts = make(map[string][]int64)

func (s *TaskService) CreateTask(task *curescan.Task) error {
	if !errors.Is(global.GVA_DB.Select("task_name", "created_by").First(&curescan.Task{}, "task_name=? AND created_by=?", task.TaskName, task.CreatedBy).Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("任务'%s'已存在, 请勿重复创建", task.TaskName)
	}

	if task.TaskPlan == common.ExecuteImmediately {
		// 普通任务创建后的状态为"创建"
		task.Status = common.Waiting
	} else if task.TaskPlan == common.ExecuteLater {
		task.Status = common.Created
	} else {
		// 定时任务创建后的状态为"停止"
		task.Status = common.TimeStopped
	}
	err := global.GVA_DB.Create(&task).Error
	if err != nil {
		return err
	}
	// 稍后执行
	if task.TaskPlan == common.ExecuteLater {
		return nil
	}

	err = global.GVA_REDIS.RPush(context.Background(), "taskQueue", task.ID).Err()
	return err
}

func (s *TaskService) UpdateTask(task *curescan.Task) error {
	var existingRecord curescan.Task
	err := global.GVA_DB.Select("id", "task_name", "created_by").Where("task_name=? AND created_by=?", task.TaskName, task.CreatedBy).First(&existingRecord).Error
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

func (s *TaskService) UpdateTaskWithTransction(tx *gorm.DB, task *curescan.Task) error {
	var existingRecord curescan.Task
	err := tx.Select("id", "task_name").Where("task_name=?", task.TaskName).First(&existingRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Save(&task).Error
		}
		return err
	}
	if existingRecord.ID != task.ID {
		return errors.New("任务名称已被占用，不允许修改")
	}
	return tx.Save(&task).Error
}

func (s *TaskService) DeleteTask(id int) error {
	return global.GVA_DB.Delete(&curescan.Task{}, id).Error
}

func (s *TaskService) GetTaskById(id int) (*curescan.Task, error) {
	var task *curescan.Task
	err := global.GVA_DB.Select("id", "task_name", "task_desc", "status", "target_ip", "policy_id", "task_plan",
		"plan_config", "created_at", "updated_at", "deleted_at", "flag", "created_by", "updated_by", "entry_id").Where("id=?", id).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("目标数据不存在")
		}
		return nil, err
	}
	policy, err := policyService.GetPolicyById(int(task.PolicyID))
	if err != nil {
		return nil, err
	}
	task.PolicyName = policy.PolicyName
	return task, nil
}

func (s *TaskService) GetTaskList(st request.SearchTask) (list interface{}, total int64, err error) {
	page := st.PageInfo
	order := st.OrderKey
	desc := st.Desc
	limit := page.PageSize
	offset := page.PageSize * (page.Page - 1)
	db := global.GVA_DB.Model(&curescan.Task{}).Select("id", "task_name", "task_desc", "status", "target_ip", "policy_id", "task_plan",
		"plan_config", "created_at", "updated_at", "deleted_at", "flag", "entry_id", "created_by")
	var tasks []*curescan.Task
	if st.TaskName != "" {
		db = db.Where("task_name LIKE ?", "%"+st.TaskName+"%")
	}
	if len(st.TaskPlan) != 0 {
		db = db.Where("task_plan in (?)", st.TaskPlan)
	}
	if st.Status != 0 {
		db = db.Where("status=?", st.Status)
	}
	if st.PolicyId != 0 {
		db = db.Where("policy_id=?", st.PolicyId)
	}

	// 数据隔离，非超管用户仅返回当前用户创建的数据
	if !st.AllData {
		if st.CreatedBy != 0 {
			db = db.Where("created_by=?", st.CreatedBy)
		}
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
	if err != nil {
		return tasks, total, err
	}
	for _, task := range tasks {
		policy, err := policyService.GetPolicyById(int(task.PolicyID))
		if err != nil {
			// todo 记录错误日志
			return tasks, total, err
		}
		task.PolicyName = policy.PolicyName
	}
	return tasks, total, err
}

// ExecuteTask 执行任务
func (s *TaskService) ExecuteTask(id int) error {
	// var wg sync.WaitGroup
	eg := errgroup.Group{}
	// 接收回调中的任务结果
	var taskResult response.TaskResult
	// 获取任务
	task, err := s.GetTaskById(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if task.Status == common.Running || task.Status == common.TimeRunning {
		return errors.New("任务正在执行中，请勿重复执行")
	}

	if task.Status == common.Success {
		task = &curescan.Task{
			CsModel: global.CsModel{
				CreatedBy: task.CreatedBy,
				UpdatedBy: task.UpdatedBy,
			},
			TaskName:   task.TaskName + "_copy_" + utils.RandomString(6),
			TaskDesc:   task.TaskDesc,
			Status:     common.Created,
			TargetIP:   task.TargetIP,
			PolicyID:   task.PolicyID,
			TaskPlan:   task.TaskPlan,
			PlanConfig: task.PlanConfig,
			Executions: 0,
			EntryID:    "",
			Flag:       task.Flag,
		}
	}

	// 得到任务关联的策略
	policy, err := policyService.GetPolicyById(int(task.PolicyID))
	if err != nil {
		return err
	}

	// 解析策略
	var onlineConfig request.OnlineConfig
	var portScanConfig request.PortScanConfig
	var jobConfig []*request.JobConfig
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

	// 构造任务参数
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
					EntryID: result.Items[0].EntryID,
				})
			}
			taskResult.PortScanList = data
			global.GVA_LOG.Info(fmt.Sprintf("任务%s端口扫描完成", task.TaskName))
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
					EntryID: result.EntryID,
				})
			}
			taskResult.OnlineCheckList = data
			global.GVA_LOG.Info(fmt.Sprintf("任务%s在线检测完成", task.TaskName))
			return nil
		},
	}
	jobs, err := s.generateJob(jobConfig, &taskResult, task)
	if err != nil {
		return err

	}
	options.Jobs = jobs
	if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		// 处理任务
		eg.Go(func() error {
			return s.processTask(task, options, &taskResult)
		})
		return eg.Wait()
	} else {
		task.Status = common.TimeRunning
		s.UpdateTask(task)
		cronName := task.TaskName
		global.GVA_Timer.AddTaskByFunc(cronName, task.PlanConfig, func() { s.processTask(task, options, &taskResult) }, task.TaskName, cron.WithSeconds())
	}

	return nil
}

// processTask 处理任务的执行流程
// 对于普通任务来说, 不需要复制任务, 但是对于定时任务每次执行需要复制一次任务
// 对于普通任务如果需要复用, 需要重新创建一条任务
func (s *TaskService) processTask(task *curescan.Task, options *types.Options, taskResult *response.TaskResult) error {
	var entry *eagleeye.EagleeyeEntry
	var err error
	if task.TaskPlan != common.ExecuteTiming {
		entry, err = global.EagleeyeEngine.NewEntry(options)
	}
	if err != nil {
		// global.GVA_LOG.Error("创建任务entry失败", zap.String("taskName", task.TaskName), zap.Error(err))
		return err
	}
	// 使用数据库事务处理整个任务流程
	if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		// task.Status = common.Running
		task.EntryID = entry.EntryID
		// err = s.UpdateTask(task)
		// if err != nil {
		// 	// global.GVA_LOG.Error("任务开始执行失败", zap.String("taskName", task.TaskName), zap.String("error", err.Error()))
		// 	return err
		// }
		return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			global.GVA_LOG.Info("任务开始执行...", zap.String("taskName", task.TaskName))
			// 执行任务的入口
			err = entry.Run(context.Background())
			if err != nil {
				return err
			}
			// result := entry.Result()
			// 任务执行成功 批量添加任务结果
			err = portScanService.BatchAddWithTransaction(tx, taskResult.PortScanList)
			if err != nil {
				return err
			}
			err = onlineCheckService.BatchAddWithTransaction(tx, taskResult.OnlineCheckList)
			if err != nil {
				return err
			}
			err = jobResultService.BatchAddWithTransaction(tx, taskResult.JobResultList)
			if err != nil {
				return err
			}
			info, err := userService.FindUserById(int(task.CreatedBy))
			if err != nil {
				return err
			}
			var indexMap = make(map[string]int, 1)
			for i, item := range options.Jobs {
				if item.Kind == common.VulnerabilityScan {
					indexMap["vuln"] = i
				}
			}
			err = s.GenerateReport(entry.Result(), info.NickName, indexMap)
			if err != nil {
				return err
			}

			// 资产添加
			assets := getAssetFromResult(taskResult, task)
			if len(assets) > 0 {
				err = assetService.BatchAddWithTransaction(tx, assets)
				if err != nil {
					return err
				}
			}
			// err = tx.Model(&curescan.Asset{}).CreateInBatches(assets, 100).Error
			if err != nil {
				return err
			}
			return nil
		})
		// if err != nil {
		// 	if errors.Is(err, eagleeye.ErrHasBeenStopped) {
		// 		// global.GVA_LOG.Error("任务终止...", zap.String("taskName", task.TaskName), zap.Error(err))
		// 		task.Status = common.Stopped
		// 		// 更新任务状态为 失败或成功
		// 		return s.UpdateTask(task)
		// 	} else {
		// 		// global.GVA_LOG.Error("任务执行失败", zap.String("taskName", task.TaskName), zap.Error(err))
		// 		return err
		// 	}
		// }

	}
	if task.TaskPlan == common.ExecuteTiming {
		newTask := &curescan.Task{
			Status:     common.Running,
			TaskName:   task.TaskName + "_" + time.Now().Format("2006-01-02 15:04:05"),
			TaskDesc:   fmt.Sprintf("%s (该任务由定时任务【%s】生成)", task.TaskDesc, task.TaskName),
			TargetIP:   task.TargetIP,
			PolicyID:   task.PolicyID,
			TaskPlan:   common.ExecuteImmediately,
			PlanConfig: "",
			// EntryID:    entry.EntryID,
			Flag:    task.Flag,
			CsModel: global.CsModel{CreatedBy: task.CreatedBy},
		}
		err = s.CreateTask(newTask)
		// err = global.GVA_DB.Create(newTask).Error
		if err != nil {
			// global.GVA_LOG.Error("任务开始执行失败", zap.String("taskName", newTask.TaskName), zap.String("error", err.Error()))
			return err
		}
	}
	return err
}

func getAssetFromResult(result *response.TaskResult, task *curescan.Task) []*curescan.Asset {
	assets := make([]*curescan.Asset, 0)
	for _, item := range result.JobResultList {
		// typeSplit := strings.Split(item.TemplateID, "_")
		if item.Kind == "1" {
			ip := strings.Split(item.Host, ":")[0]
			port, err := strconv.Atoi(item.Port)
			if err != nil {
				continue
			}
			if _, ok := ipPorts[ip]; !ok {
				ipPorts[ip] = make([]int64, 0)
			}
			ipPorts[ip] = append(ipPorts[ip], int64(port))

			if existingPort, exists := ipPortMap[ip]; !exists || portPriority[int64(port)] < portPriority[existingPort] {
				ipPortMap[ip] = int64(port)
				asset := &curescan.Asset{}
				asset.AreaName = "未知"
				asset.AssetArea = 0
				asset.AssetName = item.Name
				asset.AssetType = item.Tag1
				asset.SystemType = item.Tag2
				asset.Manufacturer = item.Tag3
				asset.AssetModel = item.Tag4
				asset.AssetIP = strings.Split(item.Host, ":")[0]
				asset.CreatedBy = task.CreatedBy
				asset.OpenPorts = ipPorts[ip]
				assets = append(assets, asset)
			}

		}
	}

	for _, item := range result.PortScanList {
		for _, asset := range assets {
			if asset.AssetIP == item.IP {
				asset.OpenPorts = item.Ports
			}
		}
		for _, port := range item.Ports {
			if assetInfo, ok := portAssetMap[port]; ok {
				ip := strings.Split(item.IP, ":")[0]
				if _, ok := ipPorts[ip]; !ok {
					ipPorts[ip] = make([]int64, 0)
				}
				ipPorts[ip] = append(ipPorts[ip], port)
				if existingPort, exists := ipPortMap[ip]; !exists || portPriority[port] < portPriority[existingPort] {
					ipPortMap[ip] = port
					asset := &curescan.Asset{}
					fmt.Println("ports", ipPorts[ip])
					asset.OpenPorts = ipPorts[ip]
					asset.AreaName = "未知"
					asset.AssetArea = 0
					asset.AssetIP = ip
					asset.AssetName = assetInfo.AssetName
					asset.AssetType = assetInfo.AssetType
					asset.AssetModel = assetInfo.AssetModel
					asset.SystemType = assetInfo.SystemType
					asset.Manufacturer = assetInfo.Manufacturer
					asset.CreatedBy = task.CreatedBy
					assets = append(assets, asset)
				}
			}
		}
	}
	return assets
}

// generateJob 生成任务， 根据任务配置生成任务
func (s *TaskService) generateJob(jobConfig []*request.JobConfig, taskResult *response.TaskResult, task *curescan.Task) ([]types.JobOptions, error) {
	jobs := make([]types.JobOptions, len(jobConfig))
	for i, job := range jobConfig {
		jobs[i].Name = job.Name
		jobs[i].Kind = job.Kind
		jobs[i].Concurrency = job.Concurrency
		jobs[i].Count = job.Count
		jobs[i].Format = job.Format
		jobs[i].Timeout = job.Timeout
		jobs[i].RateLimit = job.RateLimit
		jobs[i].ResultCallback = func(c context.Context, result *types.JobResult) error {
			for _, item := range result.Items {
				var oneRes = &curescan.JobResultItem{
					Name:             item.TemplateName,
					Kind:             result.Kind,
					TemplateID:       item.TemplateID,
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
					EntryID:          result.EntryID,
					Remediation:      item.Remediation,
					CsModel:          global.CsModel{CreatedBy: task.CreatedBy},
				}
				tagSplit := strings.Split(item.Tags, ",")
				if len(tagSplit) == 1 {
					oneRes.Tag1 = tagSplit[0]
					oneRes.Tag2 = "未知"
					oneRes.Tag3 = "未知"
					oneRes.Tag4 = "未知"
				}
				if len(tagSplit) == 2 {
					oneRes.Tag1 = tagSplit[0]
					oneRes.Tag2 = tagSplit[1]
					oneRes.Tag3 = "未知"
					oneRes.Tag4 = "未知"
				}
				if len(tagSplit) == 3 {
					oneRes.Tag1 = tagSplit[0]
					oneRes.Tag2 = tagSplit[1]
					oneRes.Tag3 = tagSplit[2]
					oneRes.Tag4 = "未知"
				}
				if len(tagSplit) == 4 {
					oneRes.Tag1 = tagSplit[0]
					oneRes.Tag2 = tagSplit[1]
					oneRes.Tag3 = tagSplit[2]
					oneRes.Tag4 = tagSplit[3]
				}

				taskResult.JobResultList = append(taskResult.JobResultList, oneRes)
			}
			global.GVA_LOG.Info(fmt.Sprintf("任务%s %s阶段执行完成", task.TaskName, result.Name))
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
					ID:       template.TemplateId,
					Original: template.TemplateContent,
				})
			}
			return rawTemplates
		}

	}
	return jobs, nil
}

func (s *TaskService) StopTask(id int) error {
	task, err := s.GetTaskById(id)
	if err != nil {
		return err
	}
	// 停止普通任务
	if task.Status == common.Waiting {
		removed := utils.RemoveValueFromList(global.GVA_REDIS, "taskQueue", strconv.Itoa(id))
		if removed {
			task.Status = common.Stopped
		} else {
			task.Status = common.Failed
		}
	} else if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		task.Status = common.Stopped
		entry := global.EagleeyeEngine.Entry(task.EntryID)
		if entry == nil {
			return errors.New("任务未开始或已结束")
		}
		if err := entry.Stop(); err != nil {
			return errors.New("任务未开始或已结束")
		}
	}

	// 停止定时任务
	if task.TaskPlan == common.ExecuteTiming {

		task.Status = common.TimeStopped
		cronName := task.TaskName
		global.GVA_Timer.StopCron(cronName)
	}
	if err := s.UpdateTask(task); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) GenerateReport(ret *types.EntryResult, reporter string, indexMap map[string]int) error {
	var err error
	ip, err := utils.GetLocalIP()
	if err != nil {
		return err
	}
	var exc []string
	for _, str := range ret.ExcludeTargets {
		if str != ip {
			exc = append(exc, str)
		}
	}
	ret.ExcludeTargets = exc

	if index, ok := indexMap["vuln"]; ok {
		err = report.Generate(
			report.WithJobIndexes(index),
			report.WithEntryResult(ret),
			report.WithCustomer(global.GVA_CONFIG.Report.Customer),
			report.WithReporter(reporter))

	} else {
		err = report.Generate(
			report.WithEntryResult(ret),
			report.WithCustomer(global.GVA_CONFIG.Report.Customer),
			report.WithReporter(reporter))
	}
	return err
}

func (s *TaskService) GetTaskStage(id int64) (*response.Stage, error) {
	task, err := s.GetTaskById(int(id))
	if err != nil {
		return nil, errors.New("目标数据不存在")
	}
	entry := global.EagleeyeEngine.Entry(task.EntryID)
	if entry == nil {
		return nil, errors.New("任务未开始或已结束")
	}
	policy, err := policyService.GetPolicyById(int(task.PolicyID))
	if err != nil {
		return nil, errors.New("目标数据的策略不存在")
	}
	modelStage := &response.Stage{}
	stage := entry.Stage()
	// 保留四位小数
	modelStage.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", stage.Percent), 64)
	var jobConfig []*response.JobConfig
	var onlineCheckConfig *response.OnlineConfig
	var portScanConfig *response.PortScanConfig
	count := 0
	err = json.Unmarshal([]byte(policy.OnlineConfig), &onlineCheckConfig)
	if err == nil && onlineCheckConfig != nil {
		count++
	}
	err = json.Unmarshal([]byte(policy.PortScanConfig), &portScanConfig)
	if err == nil && portScanConfig != nil {
		count++
	}
	err = json.Unmarshal([]byte(policy.PolicyConfig), &jobConfig)
	if err == nil && jobConfig != nil {
		count += len(jobConfig)
	}
	switch stage.Name {
	case "PortScanning":
		modelStage.Name = "端口扫描"
		modelStage.Running = 2
	case "HostDiscovery":
		modelStage.Name = "在线检测"
		modelStage.Running = 1
	case "Job":
		index, ok := stage.Entries[types.StageEntryJobIndex]
		modelStage.Running = index.(int) + (count - len(jobConfig)) + 1
		if !ok {
			return nil, errors.New("策略数据错误")
		}
		modelStage.Name = jobConfig[index.(int)].Name
	}
	modelStage.Total = count
	return modelStage, nil
}
