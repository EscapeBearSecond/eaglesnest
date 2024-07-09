package service

import (
	"47.103.136.241/goprojects/gin-vue-admin/server/service/example"
	"47.103.136.241/goprojects/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
