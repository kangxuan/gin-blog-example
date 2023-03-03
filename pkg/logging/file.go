package logging

import (
	"fmt"
	"gin-blog-example/pkg/file"
	"gin-blog-example/settings"
	"os"
	"time"
)

// getLogFilePath 获取日志目录
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", settings.AppSetting.RuntimeRootPath, settings.AppSetting.LogSavePath)
}

// getLogFileName 获取日志文件名
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s", settings.AppSetting.LogSaveName, time.Now().Format(settings.AppSetting.TimeFormat), settings.AppSetting.LogFileExt)
}

// OpenLogFile 打开日志文件
func openLogFile(fileName, filePath string) (*os.File, error) {
	// 获取当前的项目的根目录
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.GetWd err:%v", err)
	}

	// 拼装日志的目录
	src := dir + filePath
	// 检查目录的权限
	if file.CheckPermission(src) {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s\"", src)
	}

	// 检查目录是否存在并创建
	if err := file.IsNotExistedMkDir(src); err != nil {
		return nil, fmt.Errorf("file.IsNotExistedMkDir src: %s, err: %v", src, err)
	}

	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, fmt.Errorf("file.Open Error: %v", err)
	}

	return f, nil
}
