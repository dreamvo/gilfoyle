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
	Short:   "Launch background task worker",
	Example: "gilfoyle worker",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		logger.Info("Initializing worker node")
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		logger.Fatal("Not implemented yet!")
	},
}
