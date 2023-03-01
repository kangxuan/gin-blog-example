package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/" // 日志保存目录
	LogSaveName = "log"
	LogFileExt  = "log"      // 日志文件后缀
	TimeFormat  = "20060102" // 时间格式
)

// getLogFilePath 获取日志目录
func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

// getLogFileFullPath 获取日志保存路径
func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

// OpenLogFile 打开日志文件
func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err): // 判断文件不存在或目录不存在
		mkDir()
	case os.IsPermission(err): // 判断文件是不是权限不足
		log.Fatalf("Permission: %v\n", err)
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v\n", err)
	}

	return file
}

// mkDir 创建目录
func mkDir() {
	dir, _ := os.Getwd() //返回与当前目录对应的根路径名
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
