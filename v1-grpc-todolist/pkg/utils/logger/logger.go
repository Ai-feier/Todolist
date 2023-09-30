package logger

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

var LogrusObj *logrus.Logger

// init 初始化日志对象
func init() {
	if LogrusObj != nil {
		src, err := setOutFile()
		if err != nil {
			log.Println(err)
		}
		// 设置输出
		LogrusObj.Out = src
		return
	}
	// 实例化日志
	logger := logrus.New()
	src, _ := setOutFile()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
}


func setOutFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil { //nolint:wsl
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)
	if os.IsNotExist(err) { //nolint:wsl
		// 当前路径不可用, 重新创建
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	logFile := now.Format("2003-03-19") + ".log"
	// 拼接日志文件
	fileName := path.Join(logFilePath, logFile)
	if _, err = os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return src, nil
}
