package v1

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1/curescan"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1/example"
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	CurescanApiGroup curescan.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
