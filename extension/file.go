package extension

import (
	"github.com/apollogoClient/v1/env/file"
)

var fileHandler file.FileHandler

// SetFileHandler 设置备份文件处理
func SetFileHandler(inFile file.FileHandler) {
	fileHandler = inFile
}

// GetFileHandler 获取备份文件处理
func GetFileHandler() file.FileHandler {
	return fileHandler
}
