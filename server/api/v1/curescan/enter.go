package curescan

import (
	"47.103.136.241/goprojects/curescan/server/service"
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
}

var (
	areaService        = service.ServiceGroupApp.CurescanServiceGroup.AreaService
	templateService    = service.ServiceGroupApp.CurescanServiceGroup.TemplateService
	assetService       = service.ServiceGroupApp.CurescanServiceGroup.AssetService
	policyService      = service.ServiceGroupApp.CurescanServiceGroup.PolicyService
	onlineCheckService = service.ServiceGroupApp.CurescanServiceGroup.OnlineCheckService
	// portScanService    = service.ServiceGroupApp.CurescanServiceGroup.PortScanService
	taskService = service.ServiceGroupApp.CurescanServiceGroup.TaskService
	vulnService = service.ServiceGroupApp.CurescanServiceGroup.VulnService
)
