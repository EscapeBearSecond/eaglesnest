package service

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service/curescan"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service/example"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	CurescanServiceGroup curescan.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
