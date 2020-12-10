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
	rootCmd.AddCommand(workerTestCmd)
}

var workerTestCmd = &cobra.Command{
	Use:     "worker-push",
	Example: "gilfoyle worker-push",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		w, err := gilfoyle.NewWorker()
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
			logger.Fatal("Failed to initialize worker channel", zap.Error(err))
		}

		err = worker.ProduceVideoTranscodingQueue(ch, &worker.VideoTranscodingParams{
			MediaUUID:      uuid.New(),
			SourceFilePath: uuid.New().String(),
		})
		if err != nil {
			logger.Fatal("Failed to publish message", zap.Error(err))
		}

	},
}
