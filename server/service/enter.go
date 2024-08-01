package service

import (
	"47.103.136.241/goprojects/curescan/server/service/curescan"
	"47.103.136.241/goprojects/curescan/server/service/example"
	"47.103.136.241/goprojects/curescan/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	CurescanServiceGroup curescan.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
