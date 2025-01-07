package server

import (
	"bufio"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
)

// FileWriter 添加文件写入管理器结构体
type FileWriter struct {
	file   *os.File
	writer *bufio.Writer
	path   string
}

// InitFile 初始化文件
func InitFile() {
	file, err := os.Create("./File/zh_CN.json")
	if err != nil {
		log.Error("初始化文件异常 : ", err)
	}
	file.Close()
	file, err = os.Create("./File/zh_CN.lang")
	if err != nil {
		log.Error("初始化文件异常 : ", err)
	}
	file.Close()
}

// NewFileWriter 创建新的文件写入管理器
func NewFileWriter(path string) (*FileWriter, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("创建输出文件异常: %v", err)
	}

	return &FileWriter{
		file:   file,
		writer: bufio.NewWriter(file),
		path:   path,
	}, nil
}

// WriteLine 写入一行内容
func (fw *FileWriter) WriteLine(text string) error {
	_, err := fw.writer.WriteString(text + "\n")
	return err
}

// Close 关闭文件写入器
func (fw *FileWriter) Close() error {
	if err := fw.writer.Flush(); err != nil {
		return fmt.Errorf("flush文件失败: %v", err)
	}
	return fw.file.Close()
}
