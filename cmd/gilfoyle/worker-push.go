package gilfoyle

import (
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	// Register command
	rootCmd.AddCommand(workerPushCmd)
}

var workerPushCmd = &cobra.Command{
	Use: "worker-push",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		w, err := worker.New(worker.Options{
			Host:        gilfoyle.Config.Services.RabbitMQ.Host,
			Port:        gilfoyle.Config.Services.RabbitMQ.Port,
			Username:    gilfoyle.Config.Services.RabbitMQ.Username,
			Password:    gilfoyle.Config.Services.RabbitMQ.Password,
			Logger:      logger,
			Concurrency: 0,
		})
		if err != nil {
			logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
		}
		defer w.Close()

		err = w.Init()
		if err != nil {
			logger.Fatal("Failed to initialize worker queues", zap.Error(err))
		}

		ch, err := w.Client.Channel()
		if err != nil {
			logger.Fatal("channel", zap.Error(err))
		}

		for i := 0; i < 5; i++ {
			err = worker.VideoTranscodingProducer(ch, worker.VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: uuid.New().String(),
			})
			if err != nil {
				logger.Fatal("produce", zap.Error(err))
			}
		}
	},
}
