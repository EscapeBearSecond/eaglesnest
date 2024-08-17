package curescan

import (
	"47.103.136.241/goprojects/curescan/server/model/curescan/common"
	"47.103.136.241/goprojects/curescan/server/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"47.103.136.241/goprojects/curescan/server/model/curescan/response"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"47.103.136.241/goprojects/eagleeye/pkg/report"
	eagleeye "47.103.136.241/goprojects/eagleeye/pkg/sdk"
	"47.103.136.241/goprojects/eagleeye/pkg/types"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
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
	userService        = &system.UserService{}
	assetService       = &AssetService{}
)

var portAssetMap = map[int64]*curescan.Asset{
	3306: {AssetName: "MySQL", AssetType: "服务器设备", AssetModel: "MySQL", SystemType: "MySQL", Manufacturer: "MySQL"},
	3389: {AssetName: "Windows远程桌面", AssetType: "服务器设备", AssetModel: "MySQL", SystemType: "Windows", Manufacturer: "Microsoft"},
	23:   {AssetName: "telnet", AssetType: "网络设备", AssetModel: "telnet", SystemType: "telnet", Manufacturer: "telnet"},
	554:  {AssetName: "rtsp", AssetType: "视频设备", AssetModel: "rtsp", SystemType: "rtsp", Manufacturer: "rtsp"},
	// 5432: {AssetName: "PgSQL", AssetType: "服务器设备", AssetModel: "PgSQL", SystemType: "PgSQL", Manufacturer: "PgSQL"},
}

func (s *TaskService) CreateTask(task *curescan.Task) error {
	if !errors.Is(global.GVA_DB.Select("task_name").First(&curescan.Task{}, "task_name=?", task.TaskName).Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("任务'%s'已存在, 请勿重复创建", task.TaskName)
	}

	if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		// 普通任务创建后的状态为"创建"
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

	go func() {
		err := s.ExecuteTask(int(task.ID))
		if err != nil {
			return
		}
	}()
	return nil
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
	// if existingRecord.ID != task.ID {
	// 	return errors.New("任务名称已被占用，不允许修改")
	// }
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
	// if existingRecord.ID != task.ID {
	// 	return errors.New("任务名称已被占用，不允许修改")
	// }
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
		"plan_config", "created_at", "updated_at", "deleted_at", "flag", "entry_id")
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
	fmt.Println("111111111111")
	// 接收回调中的任务结果
	var taskResult response.TaskResult

	// 获取任务
	task, err := s.GetTaskById(id)

	if err != nil {
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
			return nil
		},
	}
	jobs, err := s.generateJob(jobConfig, &taskResult)
	if err != nil {
		return err

	}
	options.Jobs = jobs
	if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		// 处理任务
		go s.processTask(task, options, &taskResult)
	}
	fmt.Println("task", task.TaskName)
	if task.TaskPlan == common.ExecuteTiming {
		fmt.Println("??????????")
		task.Status = common.TimeRunning
		s.UpdateTask(task)
		cronName := task.TaskName
		global.GVA_Timer.AddTaskByFunc(cronName, task.PlanConfig, func() { s.processTask(task, options, &taskResult) }, task.TaskName, cron.WithSeconds())
		fmt.Println("list", global.GVA_Timer.FindCronList())
	}
	return nil
}

// processTask 处理任务的执行流程
// 对于普通任务来说, 不需要复制任务, 但是对于定时任务每次执行需要复制一次任务
// 对于普通任务如果需要复用, 需要重新创建一条任务
func (s *TaskService) processTask(task *curescan.Task, options *types.Options, taskResult *response.TaskResult) {
	entry, err := global.EagleeyeEngine.NewEntry(options)
	if err != nil {
		global.GVA_LOG.Error("创建任务entry失败", zap.String("taskName", task.TaskName), zap.Error(err))
		return
	}
	// 使用数据库事务处理整个任务流程
	if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		task.Status = common.Running
		task.EntryID = entry.EntryID
		err = s.UpdateTask(task)
		if err != nil {
			global.GVA_LOG.Error("任务开始执行失败", zap.String("taskName", task.TaskName), zap.String("error", err.Error()))
			return
		}
		err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
			global.GVA_LOG.Info("任务开始执行...", zap.String("taskName", task.TaskName))
			// 在Redis中记录任务和entryID的关联
			err = global.GVA_REDIS.Set(context.Background(), "task_"+strconv.Itoa(int(task.ID)), entry.EntryID, 0).Err()
			if err != nil {
				return err
			}
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
				if item.Kind == "2" {
					indexMap["vuln"] = i
				}
			}
			err = s.GenerateReport(entry.Result(), info.NickName, indexMap)
			if err != nil {
				return err
			}

			// 资产添加
			assets := getAssetFromResult(taskResult)
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
		if err != nil {
			if errors.Is(err, eagleeye.ErrHasBeenStopped) {
				global.GVA_LOG.Error("任务终止...", zap.String("taskName", task.TaskName), zap.Error(err))
				task.Status = common.Stopped
			} else {
				global.GVA_LOG.Error("任务执行失败", zap.String("taskName", task.TaskName), zap.Error(err))
				task.Status = common.Failed
			}
		} else {
			task.Status = common.Success
			global.GVA_LOG.Info("任务执行成功", zap.String("taskName", task.TaskName))
		}
		// 更新任务状态为 失败或成功
		s.UpdateTask(task)
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
			EntryID:    entry.EntryID,
			Flag:       task.Flag,
			CsModel:    global.CsModel{CreatedBy: task.CreatedBy},
		}
		err = s.CreateTask(newTask)
		// err = global.GVA_DB.Create(newTask).Error
		if err != nil {
			global.GVA_LOG.Error("任务开始执行失败", zap.String("taskName", newTask.TaskName), zap.String("error", err.Error()))
			return
		}
		s.processTask(newTask, options, taskResult)
		// err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 	// 更新任务状态为执行中
		//
		// 	global.GVA_LOG.Info("任务开始执行...", zap.String("taskName", newTask.TaskName))
		// 	// 在Redis中记录任务和entryID的关联
		// 	err = global.GVA_REDIS.Set(context.Background(), "task_"+strconv.Itoa(int(newTask.ID)), entry.EntryID, 0).Err()
		// 	if err != nil {
		// 		return err
		// 	}
		// 	// global.GVA_REDIS.Set(context.Background(), "task_"+strconv.Itoa(int(newTask.ID)), entry.EntryID, -1)
		// 	// 执行任务的入口
		// 	err = entry.Run(context.Background())
		// 	if err != nil {
		// 		return err
		// 	}
		// 	// result := entry.Result()
		// 	// fmt.Println(result)
		// 	fmt.Println("该入结果库的数据：", len(taskResult.JobResultList))
		// 	// 任务执行成功 批量添加任务结果
		// 	err = portScanService.BatchAddWithTransaction(tx, taskResult.PortScanList)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	err = onlineCheckService.BatchAddWithTransaction(tx, taskResult.OnlineCheckList)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	err = jobResultService.BatchAddWithTransaction(tx, taskResult.JobResultList)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	info, err := userService.FindUserById(int(task.CreatedBy))
		// 	if err != nil {
		// 		return err
		// 	}
		// 	var indexMap = make(map[string]int, 1)
		// 	for i, item := range options.Jobs {
		// 		if item.Kind == "2" {
		// 			indexMap["vuln"] = i
		// 		}
		// 	}
		// 	err = s.GenerateReport(entry.Result(), info.NickName, indexMap)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	// 资产添加
		// 	assets := getAssetFromResult(taskResult)
		// 	if len(assets) > 0 {
		// 		err = tx.Model(&curescan.Asset{}).CreateInBatches(assets, 100).Error
		// 		if err != nil {
		// 			return err
		// 		}
		// 	}
		// 	return nil
		// })
		// if err != nil {
		// 	if errors.Is(err, eagleeye.ErrHasBeenStopped) {
		// 		global.GVA_LOG.Error("任务终止...", zap.String("taskName", newTask.TaskName), zap.Error(err))
		// 		newTask.Status = common.Stopped
		// 	} else {
		// 		global.GVA_LOG.Error("任务执行失败", zap.String("taskName", newTask.TaskName), zap.Error(err))
		// 		newTask.Status = common.Failed
		// 	}
		// } else {
		// 	newTask.Status = common.Success
		// 	global.GVA_LOG.Info("任务执行成功", zap.String("taskName", newTask.TaskName))
		// }
		// // 更新任务状态为 失败或成功
		// s.UpdateTask(newTask)
	}

}

func getAssetFromResult(result *response.TaskResult) []*curescan.Asset {
	assets := make([]*curescan.Asset, 0)
	for _, item := range result.JobResultList {
		// typeSplit := strings.Split(item.TemplateID, "_")
		if item.Kind == "1" {
			fmt.Println("资产添加", item.Name)
			asset := &curescan.Asset{}
			asset.AreaName = "未知"
			asset.AssetArea = 0
			asset.AssetName = item.Name
			asset.AssetType = item.Tag1
			// if len(typeSplit) == 1 {
			asset.SystemType = item.Tag2
			asset.Manufacturer = item.Tag3
			asset.AssetModel = item.Tag4
			// }
			// if len(typeSplit) == 2 {
			// 	asset.SystemType = typeSplit[1]
			// 	asset.Manufacturer = "未知"
			// 	asset.AssetModel = "未知"
			// }
			// if len(typeSplit) == 3 {
			// 	asset.SystemType = typeSplit[1]
			// 	asset.Manufacturer = typeSplit[2]
			// 	asset.AssetModel = "未知"
			// }
			// if len(typeSplit) == 4 {
			// 	asset.SystemType = typeSplit[1]
			// 	asset.Manufacturer = typeSplit[2]
			// 	asset.AssetModel = typeSplit[3]
			// }
			asset.AssetIP = strings.Split(item.Host, ":")[0]
			port, _ := strconv.Atoi(item.Port)
			asset.OpenPorts = []int64{int64(port)}
			assets = append(assets, asset)
		}
	}

	for _, item := range result.PortScanList {
		for _, port := range item.Ports {
			if assetInfo, ok := portAssetMap[port]; ok {
				fmt.Println("发现端口", port, "与资产", assetInfo.AssetName, "匹配")
				asset := &curescan.Asset{
					OpenPorts:    []int64{port},
					AreaName:     "未知",
					AssetArea:    0,
					AssetIP:      fmt.Sprintf("%s:%d", item.IP, port),
					AssetName:    assetInfo.AssetName,
					AssetType:    assetInfo.AssetType,
					AssetModel:   assetInfo.AssetModel,
					SystemType:   assetInfo.SystemType,
					Manufacturer: assetInfo.Manufacturer,
				}
				assets = append(assets, asset)
			}
		}
	}
	return assets
}

// generateJob 生成任务， 根据任务配置生成任务
func (s *TaskService) generateJob(jobConfig []*request.JobConfig, taskResult *response.TaskResult) ([]types.JobOptions, error) {
	jobs := make([]types.JobOptions, len(jobConfig))
	// data := make([]*curescan.JobResultItem, 0)
	for i, job := range jobConfig {
		// dir := path.Join(global.GVA_CONFIG.AutoCode.Root, "server", "templates", job.Name)
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
				}
				tagSplit := strings.Split(item.Tags, "_")
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
	if task.TaskPlan == common.ExecuteImmediately || task.TaskPlan == common.ExecuteLater {
		if err := global.GVA_REDIS.Get(context.Background(), "task_"+strconv.Itoa(id)).Err(); err != nil {
			return errors.New("任务未执行或已结束")
		}
		entryID := global.GVA_REDIS.Get(context.Background(), "task_"+strconv.Itoa(id)).Val()
		entry := global.EagleeyeEngine.Entry(entryID)
		if entry == nil {
			return errors.New("任务未开始或已结束")
		}
		if err := entry.Stop(); err != nil {
			return errors.New("任务未开始或已结束")
		}
	}

	// 停止定时任务
	if task.TaskPlan == common.ExecuteTiming {
		// err = global.GVA_REDIS.Get(context.Background(), "cron_"+strconv.Itoa(id)).Err()
		// if err != nil {
		// 	return err
		// }
		cronName := task.TaskName
		// cronName := global.GVA_REDIS.Get(context.Background(), "cron_"+strconv.Itoa(id)).Val()
		global.GVA_Timer.StopCron(cronName)
	}
	// global.GVA_LOG.Info("任务已停止", zap.Int("id", id))
	return nil
}

func (s *TaskService) GenerateReport(ret *types.EntryResult, reporter string, indexMap map[string]int) error {
	var err error
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

// func (s *TaskService) DownloadReport(entryID string, format string) error {
// 	file, err := os.OpenFile("", os.O_RDONLY, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// }
