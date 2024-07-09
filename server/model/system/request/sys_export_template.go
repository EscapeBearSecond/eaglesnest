package request

import (
	"47.103.136.241/goprojects/gin-vue-admin/server/model/common/request"
	"47.103.136.241/goprojects/gin-vue-admin/server/model/system"
	"time"
)

type SysExportTemplateSearch struct {
	system.SysExportTemplate
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
