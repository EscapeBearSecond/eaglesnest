package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"47.103.136.241/goprojects/curescan/server/middleware"
	"github.com/gin-gonic/gin"
)

type TaskRouter struct {
}

func (t *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) {
	taskRouter := Router.Group("task").Use(middleware.OperationRecord())

	taskRouterApi := v1.ApiGroupApp.CurescanApiGroup.TaskApi
	statisticsApi := v1.ApiGroupApp.CurescanApiGroup.StatisticsApi
	{
		taskRouter.POST("", taskRouterApi.CreateTask).Use(middleware.OperationRecord())      // 创建Task
		taskRouter.GET(":id", taskRouterApi.GetTaskById)                                     // 获取单条Task消息
		taskRouter.PUT("", taskRouterApi.UpdateTask).Use(middleware.OperationRecord())       // 更新Task
		taskRouter.DELETE(":id", taskRouterApi.DeleteTask).Use(middleware.OperationRecord()) // 删除选中Task
		taskRouter.POST("list", taskRouterApi.GetTaskList)                                   // 获取Task列表
		taskRouter.GET("migrate", taskRouterApi.MigrateTable)
		taskRouter.POST("execute/:id", taskRouterApi.ExecuteTask).Use(middleware.OperationRecord())
		taskRouter.POST("stop/:id", taskRouterApi.StopTask).Use(middleware.OperationRecord())
		taskRouter.POST("report", taskRouterApi.DownloadReport).Use(middleware.OperationRecord())
		taskRouter.POST("docs", taskRouterApi.DownloadResultDocs).Use(middleware.OperationRecord())
		taskRouter.GET("stage/:id", taskRouterApi.GetTaskStage)
		taskRouter.GET("statistics", statisticsApi.GetTaskInfo)
	}
}
