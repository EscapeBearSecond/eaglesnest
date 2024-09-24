package example

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/service"

type ApiGroup struct {
	CustomerApi
	FileUploadAndDownloadApi
}

var (
	customerService              = service.ServiceGroupApp.ExampleServiceGroup.CustomerService
	fileUploadAndDownloadService = service.ServiceGroupApp.ExampleServiceGroup.FileUploadAndDownloadService
)
