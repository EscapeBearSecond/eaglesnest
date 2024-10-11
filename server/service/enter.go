package service

import (
	"github.com/EscapeBearSecond/curescan/server/service/curescan"
	"github.com/EscapeBearSecond/curescan/server/service/example"
	"github.com/EscapeBearSecond/curescan/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup   system.ServiceGroup
	ExampleServiceGroup  example.ServiceGroup
	CurescanServiceGroup curescan.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
