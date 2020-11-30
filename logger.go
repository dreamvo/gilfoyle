package gilfoyle

import (
	"net/url"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	loggerOnce sync.Once
	Logger     *zap.Logger
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func init() {
	_, err := NewLogger()
	if err != nil {
		panic(err)
	}
}

func NewLogger() (*zap.Logger, error) {
	var err error

	loggerOnce.Do(func() {
		lumberJackLogger := lumberjack.Logger{
			Filename:   "logs/server.log", //Filename is the file to write logs to
			MaxSize:    50,                //MB
			MaxBackups: 30,                //MaxBackups is the maximum number of old log files to retain.
			MaxAge:     90,                //days
			Compress:   false,
		}
		zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
			return lumberjackSink{
				Logger: &lumberJackLogger,
			}, nil
		})
		zapConfig := zap.Config{
			Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
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
			OutputPaths:      []string{"stdout", "lumberjack:logs/server.log"},
			ErrorOutputPaths: []string{"stderr"},
		}
		Logger, err = zapConfig.Build()
	})
	return Logger, err
}
