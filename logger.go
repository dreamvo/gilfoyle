package gilfoyle

import (
	"go.uber.org/zap"
	"sync"
)

var (
	once   sync.Once
	Logger *zap.Logger
)

func NewLogger() *zap.Logger {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}

		Logger = logger
	})

	return Logger
}
