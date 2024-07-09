package v1

import (
	"47.103.136.241/goprojects/gin-vue-admin/server/api/v1/example"
	"47.103.136.241/goprojects/gin-vue-admin/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
