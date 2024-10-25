package service

import (
	"github.com/EscapeBearSecond/eaglesnest/server/service/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/service/example"
	"github.com/EscapeBearSecond/eaglesnest/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup     system.ServiceGroup
	ExampleServiceGroup    example.ServiceGroup
	EaglesnestServiceGroup eaglesnest.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
