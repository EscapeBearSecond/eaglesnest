package curescan

import (
	"47.103.136.241/goprojects/curesan/server/service"
)

type ApiGroup struct {
	AreaApi
	TemplateApi
	AssetApi
}

var (
	areaService     = service.ServiceGroupApp.CurescanServiceGroup.AreaService
	templateService = service.ServiceGroupApp.CurescanServiceGroup.TemplateService
	assetService    = service.ServiceGroupApp.CurescanServiceGroup.AssetService
)
