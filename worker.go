package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/worker"
	"sync"
)

var (
	workerOnce        sync.Once
	Worker            *worker.Worker
	WorkerConcurrency uint
)

func NewWorker() (*worker.Worker, error) {
	var err error

	if WorkerConcurrency == 0 {
		WorkerConcurrency = Config.Settings.Worker.Concurrency
	}

	workerOnce.Do(func() {
		Worker, err = worker.New(worker.Options{
			Host:        Config.Services.RabbitMQ.Host,
			Port:        Config.Services.RabbitMQ.Port,
			Username:    Config.Services.RabbitMQ.Username,
			Password:    Config.Services.RabbitMQ.Password,
			Logger:      Logger,
			Concurrency: WorkerConcurrency,
		})
	})

	return Worker, err
}
