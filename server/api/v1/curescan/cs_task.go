package curescan

import (
	"encoding/json"
	"fmt"
	"strconv"

	"47.103.136.241/goprojects/curescan/server/global"
	"47.103.136.241/goprojects/curescan/server/model/common/response"
	"47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/model/curescan/request"
	"47.103.136.241/goprojects/curescan/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskApi struct {
}

func (t *TaskApi) CreateTask(c *gin.Context) {
	var createTask request.CreateTask
	err := c.ShouldBindJSON(&createTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(createTask, utils.CreateTaskVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.ValidateIP(createTask.TargetIP)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// bytes, err := json.Marshal(&createTask.PlanConfig)
	// if err != nil {
	// 	response.FailWithMessage(err.Error(), c)
	// 	return
	// }
	// var modelTask = curescan.Task{
	// 	TaskName:   createTask.TaskName,
	// 	TaskDesc:   createTask.TaskDesc,
	// 	TaskPlan:   createTask.TaskPlan,
	// 	PlanConfig: string(bytes),
	// 	PolicyID:   createTask.PolicyID,
	// 	Status:     createTask.Status,
	// 	TargetIP:   ips,
	// }
	var task = curescan.Task{
		TaskName:   createTask.TaskName,
		TaskDesc:   createTask.TaskDesc,
		TaskPlan:   createTask.TaskPlan,
		PlanConfig: createTask.PlanConfig,
		PolicyID:   createTask.PolicyID,
		TargetIP:   createTask.TargetIP,
		Flag:       createTask.Flag,
	}
	task.CreatedBy = utils.GetUserID(c)
	err = taskService.CreateTask(&task)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func (t *TaskApi) GetTaskList(c *gin.Context) {
	var searchTask request.SearchTask
	err := c.ShouldBindJSON(&searchTask)
	fmt.Println(searchTask)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(searchTask.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
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
	var modelTask = curescan.Task{
		GvaModel: global.GvaModel{
			ID: updateTask.ID,
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
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = taskService.ExecuteTask(int(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	// options := &types.Options{}
	// engine, err := eagleeye.NewEngine(eagleeye.WithDirectory("/results"))
	// if err != nil {
	// 	response.FailWithMessage(err.Error(), c)
	// 	return
	// }
	// defer engine.Close()
	// fmt.Println(options)
}

func (t *TaskApi) MigrateTable(c *gin.Context) {
	err := global.GVA_DB.AutoMigrate(&curescan.Task{})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var po = curescan.Task{
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
