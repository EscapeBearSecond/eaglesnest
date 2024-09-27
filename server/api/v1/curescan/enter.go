package curescan

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service"
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
	areaService        = service.ServiceGroupApp.CurescanServiceGroup.AreaService
	templateService    = service.ServiceGroupApp.CurescanServiceGroup.TemplateService
	assetService       = service.ServiceGroupApp.CurescanServiceGroup.AssetService
	policyService      = service.ServiceGroupApp.CurescanServiceGroup.PolicyService
	onlineCheckService = service.ServiceGroupApp.CurescanServiceGroup.OnlineCheckService
	// portScanService    = service.ServiceGroupApp.CurescanServiceGroup.PortScanService
	taskService       = service.ServiceGroupApp.CurescanServiceGroup.TaskService
	vulnService       = service.ServiceGroupApp.CurescanServiceGroup.VulnService
	resultService     = service.ServiceGroupApp.CurescanServiceGroup.JobResultService
	systemInfoService = service.ServiceGroupApp.CurescanServiceGroup.SystemInfoService
)
