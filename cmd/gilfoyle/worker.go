package gilfoyle

import (
	"log"
	"os"

	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/logging"
	"github.com/dreamvo/gilfoyle/transcoding"
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
		logger, err := logging.NewLogger(cfg.Settings.Debug, true)
		if err != nil {
			log.Fatal(err)
		}

		logger.Info("Initializing worker node")
		logger.Info("Environment", zap.Bool("debug", cfg.Settings.Debug))

		storage, err := gilfoyle.NewStorage(*cfg)
		if err != nil {
			logger.Fatal("Error initializing storage backend", zap.Error(err))
		}

		if concurrency == 0 {
			concurrency = cfg.Settings.Worker.Concurrency
		}

		if cfg.Settings.Debug {
			_ = os.Setenv("PGSSLMODE", "disable")
		}

		dbClient, err := db.NewClient(cfg.Services.DB)
		if err != nil {
			logger.Fatal("failed opening connection", zap.Error(err))
		}

		w, err := worker.New(worker.Options{
			Host:        cfg.Services.RabbitMQ.Host,
			Port:        cfg.Services.RabbitMQ.Port,
			Username:    cfg.Services.RabbitMQ.Username,
			Password:    cfg.Services.RabbitMQ.Password,
			Logger:      logger,
			Concurrency: concurrency,
			Storage:     storage,
			Database:    dbClient,
			Transcoder: transcoding.NewTranscoder(transcoding.Options{
				FFmpegBinPath: "/usr/bin/ffmpeg",
			}),
		})
		if err != nil {
			logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		}
		defer func() { _ = w.Close() }()

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
