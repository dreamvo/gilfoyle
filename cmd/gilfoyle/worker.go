package gilfoyle

import (
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var concurrency uint

func init() {
	// Register command
	rootCmd.AddCommand(workerCmd)

	workerCmd.PersistentFlags().UintVar(&concurrency, "concurrency", 3, "Number of concurrent messages this worker node can handle at the same time. Constraints: (1 <= n <= 1000). Concurrency N will produce N goroutines for each queue.")
}

var workerCmd = &cobra.Command{
	Use:     "worker",
	Short:   "Launch a background task worker node",
	Long:    "Multiple worker nodes represent a worker pool. We usually recommend to launch a minimum of 3 worker nodes to ensure automatic fail over and high availability.",
	Example: "gilfoyle worker",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		logger.Info("Initializing worker node")
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		storage, err := gilfoyle.NewStorage(config.StorageClass(gilfoyle.Config.Storage.Class))
		if err != nil {
			logger.Fatal("Error initializing storage backend", zap.Error(err))
		}

		if concurrency == 0 {
			concurrency = gilfoyle.Config.Settings.Worker.Concurrency
		}

		w, err := worker.New(worker.Options{
			Host:        gilfoyle.Config.Services.RabbitMQ.Host,
			Port:        gilfoyle.Config.Services.RabbitMQ.Port,
			Username:    gilfoyle.Config.Services.RabbitMQ.Username,
			Password:    gilfoyle.Config.Services.RabbitMQ.Password,
			Logger:      gilfoyle.Logger,
			Concurrency: concurrency,
			Storage:     storage,
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

		logger.Info("Worker is now ready to handle incoming messages")

		err = w.Consume()
		if err != nil {
			logger.Fatal("Failed to start consuming worker queues", zap.Error(err))
		}

		<-forever
	},
}
