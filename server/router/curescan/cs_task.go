package curescan

import (
	v1 "47.103.136.241/goprojects/curesan/server/api/v1"
	"github.com/gin-gonic/gin"
)

type TaskRouter struct {
}

func (t *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) {
	taskRouter := Router.Group("task")

	taskRouterApi := v1.ApiGroupApp.CurescanApiGroup.TaskApi
	{
		taskRouter.POST("", taskRouterApi.CreateTask) // 创建Area
		// areaRouter.POST("deleteArea", areaRouterApi.DeleteApi)               // 删除Area
		taskRouter.GET(":id", taskRouterApi.GetTaskById)   // 获取单条Area消息
		taskRouter.PUT("", taskRouterApi.UpdateTask)       // 更新Area
		taskRouter.DELETE(":id", taskRouterApi.DeleteTask) // 删除选中Area

		// apiRouterWithoutRecord.POST("getAllAreas", areaRouterApi.GetAllApis) // 获取所有Area
		taskRouter.GET("", taskRouterApi.GetTaskList) // 获取Area列表
		taskRouter.GET("migrate", taskRouterApi.MigrateTable)

	}
}
