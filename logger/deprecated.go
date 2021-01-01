// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var (
	logger      *zap.Logger
	loggerSugar *zap.SugaredLogger
)

// SetLoggerWriter
func SetLoggerWriter(path string) io.Writer {
	return &lumberjack.Logger{
		Filename:   path,
		MaxSize:    20,    //单个日志文件最大MaxSize*M大小 // megabytes
		MaxAge:     30,    //days
		MaxBackups: 50,    //备份数量
		Compress:   false, //不压缩
		LocalTime:  true,  //备份名采用本地时间
	}
}

// InitLogger 初始化日志组件
func InitLogger(path, level string) {
	log.Println(path, level)
	if path == "" || level == "" {
		logger = NewZap("DEBUG", zapcore.NewConsoleEncoder, os.Stdout)
	} else {
		logger = NewZap(level, zapcore.NewJSONEncoder, SetLoggerWriter(path))
	}
	loggerSugar = logger.Sugar()
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func Infof(format string, v ...interface{}) {
	if loggerSugar != nil {
		loggerSugar.Infof(format, v...)
	}
}

func Info(message string) {
	if logger != nil {
		logger.Info(message)
	}
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func Debugf(format string, v ...interface{}) {
	if loggerSugar != nil {
		loggerSugar.Debugf(format, v...)
	}
}

func Debug(message string) {
	if logger != nil {
		logger.Debug(message)
	}
}

func Warnf(format string, v ...interface{}) {
	if loggerSugar != nil {
		loggerSugar.Warnf(format, v...)
	}
}

func Warn(message string) {
	if logger != nil {
		logger.Warn(message)
	}
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func Errorf(format string, v ...interface{}) {
	if loggerSugar != nil {
		loggerSugar.Errorf(format, v...)
	}
}

func Error(message string) {
	if logger != nil {
		logger.Error(message)
	}
}

func Fatalf(format string, v ...interface{}) {
	if loggerSugar != nil {
		loggerSugar.Fatalf(format, v...)
	}
}

func Fatal(message string) {
	if logger != nil {
		logger.Fatal(message)
	}
}
