package gilfoyle

import (
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/dashboard"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var endpoint string

func init() {
	// Register command
	rootCmd.AddCommand(dashboardCmd)

	// Register flags
	dashboardCmd.PersistentFlags().StringVar(&addr, "addr", "", "Interface binding for the web server")
	dashboardCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 3000, "HTTP port")
	dashboardCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "http://localhost:3001", "Endpoint")
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch a web server to interact with your Gilfoyle instance",
	Long:  "Internal server requests can be made from the client through a proxy at /api/proxy.",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		logger.Info(fmt.Sprintf("Initializing Dashboard web server on port %v", httpPort))
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		if gilfoyle.Config.Settings.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		server := dashboard.NewServer(logger, endpoint)

		// Launch web server
		if err := server.Listen(fmt.Sprintf("%s:%d", addr, httpPort)); err != nil {
			logger.Fatal("error while launching web server", zap.Error(err))
		}
	},
}
