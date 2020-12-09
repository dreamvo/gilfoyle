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
	Short:   "Launch background task worker",
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

		w.Init()
	},
}
