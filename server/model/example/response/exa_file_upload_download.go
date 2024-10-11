package response

import "github.com/EscapeBearSecond/curescan/server/model/example"

type ExaFileResponse struct {
	File example.ExaFileUploadAndDownload `json:"file"`
}
