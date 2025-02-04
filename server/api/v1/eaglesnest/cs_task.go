package eaglesnest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/common"
	"github.com/EscapeBearSecond/eaglesnest/server/service/system"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/EscapeBearSecond/eaglesnest/server/global"
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/response"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest/request"
	"github.com/EscapeBearSecond/eaglesnest/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskApi struct {
}

func (t *TaskApi) CreateTask(c *gin.Context) {
	var createTask request.CreateTask
	err := c.ShouldBindJSON(&createTask)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.String("uri", c.Request.URL.Path), zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	err = utils.Verify(createTask, utils.CreateTaskVerify)
	if err != nil {
		response.FailWithMessage("请求数据不正确!", c)
		return
	}
	if len(createTask.AreaIDArray) == 0 && len(createTask.TargetIP) == 0 {
		response.FailWithMessage("请选择扫描目标", c)
		return
	}
	err = utils.ValidateIP(createTask.TargetIP)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if createTask.TaskPlan == common.ExecuteTiming && !utils.IsValidCron(createTask.PlanConfig) {
		response.FailWithMessage("非法的cron达式", c)
		return
	}

	var task = &eaglesnest.Task{
		TaskName:    createTask.TaskName,
		TaskDesc:    createTask.TaskDesc,
		TaskPlan:    createTask.TaskPlan,
		PlanConfig:  createTask.PlanConfig,
		PolicyID:    createTask.PolicyID,
		TargetIP:    createTask.TargetIP,
		Flag:        createTask.Flag,
		AreaIDArray: createTask.AreaIDArray,
		CsModel:     global.CsModel{CreatedBy: utils.GetUserID(c)},
	}
	task.CreatedBy = utils.GetUserID(c)
	if len(task.AreaIDArray) != 0 {
		for _, areaId := range task.AreaIDArray {
			area, err := areaService.GetAreaById(int(areaId))
			if err != nil {
				global.GVA_LOG.Error("创建失败!", zap.String("uri", c.Request.URL.Path), zap.Error(err))
				response.FailWithMessage("创建失败!", c)
				return
			}
			task.TargetIP = append(task.TargetIP, area.AreaIP...)
		}
	}
	err = taskService.CreateTask(task)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.String("uri", c.Request.URL.Path), zap.Error(err))
		response.FailWithMessage("创建失败!", c)
		return
	}
	response.Ok(c)
}

func (t *TaskApi) GetTaskList(c *gin.Context) {
	var searchTask request.SearchTask
	err := c.ShouldBindJSON(&searchTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(searchTask.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	searchTask.CreatedBy = utils.GetUserID(c)
	searchTask.AllData = system.HasAllDataAuthority(c)
	list, total, err := taskService.GetTaskList(searchTask)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     searchTask.Page,
		PageSize: searchTask.PageSize,
	}, "获取成功", c)
}

func (t *TaskApi) UpdateTask(c *gin.Context) {
	var updateTask request.UpdateTask
	err := c.ShouldBindJSON(&updateTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(updateTask.CreateTask, utils.CreateTaskVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var ips = updateTask.TargetIP
	err = utils.ValidateIP(ips)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bytes, err := json.Marshal(&updateTask.PlanConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var modelTask = eaglesnest.Task{
		GvaModel: global.GvaModel{
			ID: updateTask.ID,
		},
		CsModel: global.CsModel{
			UpdatedBy: utils.GetUserID(c),
			CreatedBy: utils.GetUserID(c),
		},
		TaskName:   updateTask.TaskName,
		TaskDesc:   updateTask.TaskDesc,
		TaskPlan:   updateTask.TaskPlan,
		PlanConfig: string(bytes),
		PolicyID:   updateTask.PolicyID,
		Status:     updateTask.Status,
		TargetIP:   ips,
	}
	err = taskService.UpdateTask(&modelTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (t *TaskApi) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	task, err := taskService.GetTaskById(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if task.Status == common.Running || task.Status == common.TimeRunning || task.Status == common.Waiting {
		response.FailWithMessage("任务处于执行中或队列中，不允许删除", c)
		return
	}
	err = taskService.DeleteTask(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (t *TaskApi) GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	task, err := taskService.GetTaskById(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(task, c)
}

func (t *TaskApi) ExecuteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("参数有误", c)
		return
	}
	task, err := taskService.GetTaskById(int(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if task.Status != common.Success {

		err = global.GVA_REDIS.RPush(context.Background(), "taskQueue", id).Err()
		if err != nil {
			response.FailWithMessage("加入执行队列失败", c)
			return
		}
		task.Status = common.Waiting
		err = taskService.UpdateTask(task)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.Ok(c)
	} else {
		newTask := &eaglesnest.Task{
			CsModel: global.CsModel{
				CreatedBy: task.CreatedBy,
				UpdatedBy: task.UpdatedBy,
			},
			TaskName:    strings.Split(task.TaskName, "_")[0] + "_" + time.Now().Format("2006-01-02 15:04:05"),
			TaskDesc:    task.TaskDesc,
			Status:      common.Created,
			TargetIP:    task.TargetIP,
			PolicyID:    task.PolicyID,
			TaskPlan:    common.ExecuteImmediately,
			PlanConfig:  task.PlanConfig,
			Executions:  0,
			EntryID:     "",
			Flag:        task.Flag,
			AreaIDArray: task.AreaIDArray,
		}
		err := taskService.CreateTask(newTask)
		if err != nil {
			response.FailWithMessage("执行失败", c)
			return
		}
		response.Ok(c)
	}
}

func (t *TaskApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&eaglesnest.Task{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var po = eaglesnest.Task{
		TaskName:   "",
		TaskDesc:   "",
		Status:     0,
		TargetIP:   nil,
		PolicyID:   0,
		TaskPlan:   0,
		PlanConfig: "",
		Executions: 0,
	}
	b, _ := json.Marshal(po)
	fmt.Println(string(b))
	response.Ok(c)
}

func (t *TaskApi) StopTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = taskService.StopTask(int(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (t *TaskApi) DownloadReport(c *gin.Context) {
	entryID := c.PostForm("entryId")
	format := c.PostForm("format")
	if entryID == "" || format == "" {
		response.FailWithMessage("参数有误", c)
		return
	}
	dir := filepath.Join(global.GVA_CONFIG.AutoCode.Root, "server", "reports")
	fileName := "report_" + entryID + "." + format
	fullPath := filepath.Join(dir, fileName)
	if utils.FileExists(fullPath) {
		file, err := os.Open(fullPath)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		defer file.Close()

		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		switch format {
		case "docx":
			c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
		case "pdf":
			c.Writer.Header().Set("Content-Type", "application/pdf")
		case "html":
			c.Writer.Header().Set("Content-Type", "text/html")
		default:
			c.Writer.Header().Set("Content-Type", "application/octet-stream")
		}

		if _, err := io.Copy(c.Writer, bufio.NewReader(file)); err != nil {
			global.GVA_LOG.Error("下载文件失败", zap.Error(err))
			c.Writer.Header().Set("Content-Type", "application/octet-stream")
			response.FailWithMessage("下载文件失败", c)
			return
		}
		// response.Ok(c)
	} else {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		response.FailWithMessage("当前任务暂无数据", c)
		return
	}
}

func (t *TaskApi) DownloadResultDocs(c *gin.Context) {
	entryID := c.PostForm("entryId")
	if entryID == "" {
		response.FailWithMessage("参数有误", c)
		return
	}
	dir := filepath.Join(global.GVA_CONFIG.AutoCode.Root, "server", "results", entryID)
	existed, _ := utils.PathExists(dir)
	if !existed {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		response.FailWithMessage("当前任务暂无数据", c)
		return
	}
	buf, err := utils.CreateZipFromDir(dir)
	if err != nil {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		response.FailWithMessage("获取结果文件失败", c)
		return
	}
	// 设置 Content-Disposition 头部以指示浏览器下载文件
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", entryID))
	c.DataFromReader(http.StatusOK, int64(buf.Len()), "application/zip", buf, nil)
}

func (t *TaskApi) GetTaskStage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessage("参数有误", c)
		return
	}
	if id == 0 {
		response.FailWithMessage("参数有误", c)
		return
	}
	stage, err := taskService.GetTaskStage(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(stage, c)
}
