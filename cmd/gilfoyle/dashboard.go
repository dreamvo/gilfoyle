package gilfoyle

import (
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/gin-gonic/gin"
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
		logger := gilfoyle.Logger

		logger.Info(fmt.Sprintf("Initializing Dashboard web server on port %v", httpPort))
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		if !gilfoyle.Config.Settings.Debug {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}

		router := gin.New()

		router.Use(func(ctx *gin.Context) {
			ctx.Next()

			path := ctx.Request.URL.Path
			raw := ctx.Request.URL.RawQuery
			errorMsg := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			log := logger.With(
				zap.String("Method", ctx.Request.Method),
				zap.String("Path", path),
				zap.Int("StatusCode", ctx.Writer.Status()),
				zap.Int("BodySize", ctx.Writer.Size()),
				zap.String("UserAgent", ctx.Request.UserAgent()),
			)

			if errorMsg != "" {
				log.Error("Incoming HTTP Request",
					zap.String("ErrorMessage", errorMsg),
				)
				return
			}

			log.Info("Incoming HTTP Request")
		})

		router.Use(gin.Recovery())

		// TODO(sundowndev): fix static assets serving
		router.Static("/", "../../dashboard/dist/")

		// Launch web server
		if err := router.Run(fmt.Sprintf(":%d", httpPort)); err != nil {
			logger.Fatal("error while launching web server", zap.Error(err))
		}
	},
}
