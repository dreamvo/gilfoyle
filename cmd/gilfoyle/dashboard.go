package gilfoyle

import (
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	// Register command
	rootCmd.AddCommand(dashboardCmd)

	// Register flags
	dashboardCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 3000, "HTTP port")
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch the web UI to interact with your Gilfoyle instance",
	Run: func(cmd *cobra.Command, args []string) {
		gilfoyle.Logger.Info(fmt.Sprintf("Initializing Dashboard web server on port %v", httpPort))
		gilfoyle.Logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))
	},
}
