package worker

import (
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
)

const (
	VideoTranscodingQueue    string = "VideoTranscoding"
	ThumbnailGenerationQueue string = "ThumbnailGeneration"
	PreviewGenerationQueue   string = "PreviewGeneration"
)

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	Handler    func(*Worker, <-chan amqp.Delivery)
}

var queues = []Queue{
	{
		Name:       VideoTranscodingQueue,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
		Handler:    videoTranscodingQueueConsumer,
	},
	{
		Name:       ThumbnailGenerationQueue,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
		Handler:    func(*Worker, <-chan amqp.Delivery) {},
	},
	{
		Name:       PreviewGenerationQueue,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
		Handler:    func(*Worker, <-chan amqp.Delivery) {},
	},
}

type Options struct {
	Host     string
	Port     int16
	Username string
	Password string
	Logger   *zap.Logger
}

type Worker struct {
	m      *sync.RWMutex
	Queues map[string]amqp.Queue
	Logger *zap.Logger
	Client *amqp.Connection
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
		Queues: map[string]amqp.Queue{},
		Client: conn,
		Logger: opts.Logger,
		m:      &sync.RWMutex{},
	}, nil
}

func (w *Worker) Init() error {
	ch, err := w.Client.Channel()
	if err != nil {
		return err
	}

	for _, q := range queues {
		queue, err := ch.QueueDeclare(
			q.Name,       // name
			q.Durable,    // durable
			q.AutoDelete, // delete when unused
			q.Exclusive,  // exclusive
			q.NoWait,     // no-wait
			q.Args,       // arguments
		)
		if err != nil {
			return err
		}

		w.Queues[q.Name] = queue
	}

	return nil
}

func (w *Worker) Consume() {
	ch, err := w.Client.Channel()
	if err != nil {
		w.Logger.Fatal("Error creating channel", zap.Error(err))
		return
	}

	for _, q := range queues {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			map[string]interface{}{},
		)
		if err != nil {
			w.Logger.Sugar().Fatalf("Error consuming %s queue: %e", q.Name, err)
			return
		}

		go q.Handler(w, msgs)
	}
}

func (w *Worker) Close() error {
	err := w.Client.Close()
	if err != nil {
		return err
	}

	return nil
}
