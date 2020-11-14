package gilfoyle

import (
	"go.uber.org/zap"
	"sync"
)

var (
	loggerOnce sync.Once
	Logger     *zap.Logger
)

func init() {
	_, err := NewLogger()
	if err != nil {
		panic(err)
	}
}

func NewLogger() (*zap.Logger, error) {
	var err error

	loggerOnce.Do(func() {
		Logger, err = zap.NewProduction()
	})

	return Logger, err
}
