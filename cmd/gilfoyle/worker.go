package gilfoyle

import (
	"github.com/dreamvo/gilfoyle"
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
	Long:    "Multiple worker nodes represent a worker pool. We usually recommend to launch a minimum of 2 worker nodes to ensure fail over.",
	Example: "gilfoyle worker",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		logger.Info("Initializing worker node")
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		w, err := gilfoyle.NewWorker()
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

		logger.Info("Worker is now ready to handle incoming messages")

		<-forever
	},
}
