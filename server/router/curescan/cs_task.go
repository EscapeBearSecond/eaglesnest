package curescan

import (
	v1 "47.103.136.241/goprojects/curescan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type TaskRouter struct {
}

func (t *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) {
	taskRouter := Router.Group("task")

	taskRouterApi := v1.ApiGroupApp.CurescanApiGroup.TaskApi
	{
		taskRouter.POST("", taskRouterApi.CreateTask)      // 创建Task
		taskRouter.GET(":id", taskRouterApi.GetTaskById)   // 获取单条Task消息
		taskRouter.PUT("", taskRouterApi.UpdateTask)       // 更新Task
		taskRouter.DELETE(":id", taskRouterApi.DeleteTask) // 删除选中Task
		taskRouter.POST("list", taskRouterApi.GetTaskList) // 获取Task列表
		taskRouter.GET("migrate", taskRouterApi.MigrateTable)
		taskRouter.POST("execute/:id", taskRouterApi.ExecuteTask)
		taskRouter.POST("stop/:id", taskRouterApi.StopTask)
		taskRouter.POST("report", taskRouterApi.DownloadReport)
		taskRouter.POST("docs", taskRouterApi.DownloadResultDocs)
	}
}
