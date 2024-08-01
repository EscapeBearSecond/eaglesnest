package response

import "47.103.136.241/goprojects/curescan/server/model/example"

type ExaFileResponse struct {
	File example.ExaFileUploadAndDownload `json:"file"`
}
