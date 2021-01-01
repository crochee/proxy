// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *Logger

type Logger struct {
	level       zapcore.Level
	path        string
	logger      *zap.Logger
	loggerSugar *zap.SugaredLogger
}

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
func InitLogger(opts ...Option) {
	logger = &Logger{}
	for _, opt := range opts {
		opt.Apply(logger)
	}
	if logger.path == "" {
		logger.logger = NewZap(logger.level, zapcore.NewConsoleEncoder, os.Stdout)
	} else {
		logger.logger = NewZap(logger.level, zapcore.NewConsoleEncoder, SetLoggerWriter(logger.path))
	}
	logger.loggerSugar = logger.logger.Sugar()
}

// Infof 打印Info信息
//
// @param: format 格式信息
// @param: v 参数信息
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Infof(format, v...)
	}
}

func Info(message string) {
	logger.Info(message)
}

func (l *Logger) Info(message string) {
	if l != nil {
		l.logger.Info(message)
	}
}

// Debugf 打印Debug信息
//
// @param: format 格式信息
// @param: v 参数信息
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Debugf(format, v...)
	}
}

func Debug(message string) {
	logger.Debug(message)
}

func (l *Logger) Debug(message string) {
	if l != nil {
		l.logger.Debug(message)
	}
}

func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Warnf(format, v...)
	}
}

func Warn(message string) {
	logger.Warn(message)
}

func (l *Logger) Warn(message string) {
	if l != nil {
		l.logger.Warn(message)
	}
}

// Errorf 打印Error信息
//
// @param: format 格式信息
// @param: v 参数信息
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Errorf(format, v...)
	}
}

func Error(message string) {
	logger.Error(message)
}

func (l *Logger) Error(message string) {
	if l != nil {
		l.logger.Error(message)
	}
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l != nil {
		l.loggerSugar.Fatalf(format, v...)
	}
}

func Fatal(message string) {
	logger.Fatal(message)
}

func (l *Logger) Fatal(message string) {
	if l != nil {
		l.logger.Fatal(message)
	}
}
