package curescan

import (
	"47.103.136.241/goprojects/curesan/server/service"
)

type ApiGroup struct {
	AreaApi
	TemplateApi
}

var (
	areaService = service.ServiceGroupApp.CurescanServiceGroup.AreaService
)
