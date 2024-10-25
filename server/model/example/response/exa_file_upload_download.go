package response

import "github.com/EscapeBearSecond/eaglesnest/server/model/example"

type ExaFileResponse struct {
	File example.ExaFileUploadAndDownload `json:"file"`
}
