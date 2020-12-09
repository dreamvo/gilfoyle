package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/worker"
	"sync"
)

var (
	workerOnce sync.Once
	Worker     *worker.Worker
)

func NewWorker() (*worker.Worker, error) {
	var err error

	workerOnce.Do(func() {
		Worker, err = worker.New(worker.Options{
			Host:     Config.Services.RabbitMQ.Host,
			Port:     Config.Services.RabbitMQ.Port,
			Username: Config.Services.RabbitMQ.Username,
			Password: Config.Services.RabbitMQ.Password,
			Logger:   Logger,
		})
	})

	return Worker, err
}
