package gilfoyle

import (
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	// Register command
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:     "worker",
	Short:   "Launch a background task worker node",
	Long:    "Multiple worker nodes compose a worker pool. We usually recommend to launch a minimum of 2 worker nodes to ensure fail over.",
	Example: "gilfoyle worker",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		logger.Info("Initializing worker node")
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		w, err := worker.New(worker.Options{
			Host:     gilfoyle.Config.Services.RabbitMQ.Host,
			Port:     gilfoyle.Config.Services.RabbitMQ.Port,
			Username: gilfoyle.Config.Services.RabbitMQ.Username,
			Password: gilfoyle.Config.Services.RabbitMQ.Password,
			Logger:   logger,
		})
		if err != nil {
			logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		}
		defer w.Close()

		forever := make(chan bool)

		err = w.Init()
		if err != nil {
			logger.Fatal("Failed to initialize worker queues", zap.Error(err))
		}

		w.Consume()

		ch, err := w.Client.Channel()
		if err != nil {
			logger.Fatal("Failed to create message queue channel", zap.Error(err))
		}
		defer ch.Close()

		//time.Sleep(2 * time.Second)
		//err = ch.Publish("", worker.VideoTranscodingQueue, false, false, amqp.Publishing{
		//	DeliveryMode: amqp.Persistent,
		//	ContentType:  "text/plain",
		//	Body:         []byte("hello!!!"),
		//})
		//if err != nil {
		//	logger.Error("Failed to publish a message", zap.Error(err))
		//}

		logger.Info("Worker is now ready to handle incoming jobs")

		<-forever
	},
}
