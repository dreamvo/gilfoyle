package worker

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
//videoTranscodingQueue string = "videoTranscoding"
//thumbnailGenerationQueue string = "thumbnailGeneration"
//previewGenerationQueue string = "previewGeneration"
)

type Options struct {
	Host     string
	Port     int16
	Username string
	Password string
	Logger   *zap.Logger
}

type Worker struct {
	client *amqp.Connection
	logger *zap.Logger
}

func New(opts Options) (*Worker, error) {
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
	))
	if err != nil {
		return nil, err
	}

	return &Worker{
		client: conn,
		logger: opts.Logger,
	}, nil
}

func (w *Worker) Init() {
	w.logger.Info("test")
}

func (w *Worker) Close() error {
	return w.client.Close()
}
