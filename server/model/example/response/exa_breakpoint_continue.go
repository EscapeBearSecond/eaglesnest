package response

import "github.com/EscapeBearSecond/eaglesnest/server/model/example"

type FilePathResponse struct {
	FilePath string `json:"filePath"`
}

type FileResponse struct {
	File example.ExaFile `json:"file"`
}
