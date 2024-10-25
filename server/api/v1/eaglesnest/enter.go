package eaglesnest

import (
	"github.com/EscapeBearSecond/eaglesnest/server/service"
)

type ApiGroup struct {
	AreaApi
	TemplateApi
	AssetApi
	PolicyApi
	OnlineCheckApi
	PortScanApi
	TaskApi
	VulnApi
	StatisticsApi
	SystemInfoApi
}

var (
	areaService        = service.ServiceGroupApp.EaglesnestServiceGroup.AreaService
	templateService    = service.ServiceGroupApp.EaglesnestServiceGroup.TemplateService
	assetService       = service.ServiceGroupApp.EaglesnestServiceGroup.AssetService
	policyService      = service.ServiceGroupApp.EaglesnestServiceGroup.PolicyService
	onlineCheckService = service.ServiceGroupApp.EaglesnestServiceGroup.OnlineCheckService
	// portScanService    = service.ServiceGroupApp.EaglesnestServiceGroup.PortScanService
	taskService       = service.ServiceGroupApp.EaglesnestServiceGroup.TaskService
	vulnService       = service.ServiceGroupApp.EaglesnestServiceGroup.VulnService
	resultService     = service.ServiceGroupApp.EaglesnestServiceGroup.JobResultService
	systemInfoService = service.ServiceGroupApp.EaglesnestServiceGroup.SystemInfoService
)
