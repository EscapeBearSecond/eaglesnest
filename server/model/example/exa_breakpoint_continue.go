package example

import (
	"codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/global"
)

// file struct, 文件结构体
type ExaFile struct {
	global.GvaModel
	FileName     string
	FileMd5      string
	FilePath     string
	ExaFileChunk []ExaFileChunk
	ChunkTotal   int
	IsFinish     bool
}

// file chunk struct, 切片结构体
type ExaFileChunk struct {
	global.GvaModel
	ExaFileID       uint
	FileChunkNumber int
	FileChunkPath   string
}
