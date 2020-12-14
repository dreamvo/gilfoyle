package logging

import "go.uber.org/zap"

type ILogger interface {
	Error(msg string, field ...zap.Field)
	Info(msg string, field ...zap.Field)
	Fatal(msg string, field ...zap.Field)
}
