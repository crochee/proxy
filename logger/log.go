// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package logger

import (
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Encoder func(zapcore.EncoderConfig) zapcore.Encoder

func NewZap(level string, encoderFunc Encoder, w io.Writer, fields ...zap.Field) *zap.Logger {
	core := zapcore.NewCore(
		encoderFunc(newEncoderConfig()),
		zap.CombineWriteSyncers(zapcore.AddSync(w)),
		newLevel(level),
	)
	//大于error增加堆栈信息
	return zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.DPanicLevel))
}

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

type Level string

const (
	DebugLevel  Level = "DEBUG"
	InfoLevel   Level = "INFO"
	WarnLevel   Level = "WARN"
	ErrorLevel  Level = "ERROR"
	DpanicLevel Level = "DPANIC"
	PanicLevel  Level = "PANIC"
	FatalLevel  Level = "FATAL"
)

func newLevel(level string) zapcore.Level {
	if coreLevel, ok := map[Level]zapcore.Level{
		DebugLevel:  zap.DebugLevel,
		InfoLevel:   zap.InfoLevel,
		WarnLevel:   zap.WarnLevel,
		ErrorLevel:  zap.ErrorLevel,
		DpanicLevel: zap.DPanicLevel,
		PanicLevel:  zap.PanicLevel,
		FatalLevel:  zap.FatalLevel,
	}[Level(strings.ToUpper(level))]; ok {
		return coreLevel
	}
	return zap.InfoLevel
}
