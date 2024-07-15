package curescan

import (
	"47.103.136.241/goprojects/curesan/server/service"
)

type ApiGroup struct {
	AreaApi
	TemplateApi
	AssetApi
	PolicyApi
	OnlineCheckApi
	PortScanApi
}

var (
	areaService        = service.ServiceGroupApp.CurescanServiceGroup.AreaService
	templateService    = service.ServiceGroupApp.CurescanServiceGroup.TemplateService
	assetService       = service.ServiceGroupApp.CurescanServiceGroup.AssetService
	policyService      = service.ServiceGroupApp.CurescanServiceGroup.PolicyService
	onlineCheckService = service.ServiceGroupApp.CurescanServiceGroup.OnlineCheckService
	portScanService    = service.ServiceGroupApp.CurescanServiceGroup.PortScanService
)
