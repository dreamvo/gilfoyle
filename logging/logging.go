package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/url"
)

const (
	LogFileName string = "server.log"
)

type ILogger interface {
	Error(msg string, field ...zap.Field)
	Info(msg string, field ...zap.Field)
	Fatal(msg string, field ...zap.Field)
	With(field ...zap.Field) *zap.Logger
	//Sugar() *zap.SugaredLogger
	//Named(s string) *zap.Logger
	//WithOptions(opts ...zap.Option) *zap.Logger
	//Debug(msg string, fields ...zap.Field)
	//Named(s string) *zap.Logger
	//Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry
	//DPanic(msg string, fields ...zap.Field)
	//Panic(msg string, fields ...zap.Field)
	//Core() zapcore.Core
	Sync() error
}

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func NewLumberjackSink(*url.URL) (zap.Sink, error) {
	lumberJackLogger := lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%s", LogFileName), // Filename is the file to write logs to.
		MaxSize:    50,                                  // Max file size in Mb.
		MaxBackups: 30,                                  // MaxBackups is the maximum number of old log files to retain.
		MaxAge:     90,                                  // Max file age in days.
		Compress:   false,
	}
	return lumberjackSink{
		Logger: &lumberJackLogger,
	}, nil
}

func NewLogger(debug bool, saveOnDisk bool) (ILogger, error) {
	if saveOnDisk {
		err := zap.RegisterSink("lumberjack", NewLumberjackSink)
		if err != nil {
			return nil, err
		}
	}

	var logLevel zapcore.Level
	if debug {
		logLevel = zapcore.DebugLevel
	} else {
		logLevel = zapcore.InfoLevel
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if saveOnDisk {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, fmt.Sprintf("lumberjack:logs/%s", LogFileName))
	}

	return zapConfig.Build()
}
