package gilfoyle

import (
	"context"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent/migrate"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var httpPort int

func init() {
	// Register command
	rootCmd.AddCommand(serveCmd)

	// Register flags
	serveCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 3000, "HTTP port")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch the REST API",
	Run: func(cmd *cobra.Command, args []string) {
		logger := gilfoyle.Logger

		logger.Info("Initializing API server")
		logger.Info("Environment", zap.Bool("debug", gilfoyle.Config.Settings.Debug))

		if !gilfoyle.Config.Settings.Debug {
			gin.SetMode(gin.ReleaseMode)
		} else {
			_ = os.Setenv("PGSSLMODE", "disable")
			gin.SetMode(gin.DebugMode)
		}

		err := db.InitClient(&gilfoyle.Config)
		if err != nil {
			logger.Fatal("failed opening connection: %v", zap.Error(err))
		}
		defer db.Client.Close()

		// run the auto migration tool.
		if err := db.Client.Schema.Create(
			context.Background(),
			migrate.WithDropIndex(true),
			migrate.WithDropColumn(true),
		); err != nil {
			logger.Fatal("failed creating schema resources", zap.Error(err))
		}

		logger.Info("Successfully executed database auto migration")

		router := gin.New()

		router.Use(gin.Recovery())

		router.Use(func(ctx *gin.Context) {
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

		api.RegisterRoutes(router)

		// Launch web server
		if err := router.Run(fmt.Sprintf(":%d", httpPort)); err != nil {
			logger.Fatal("error while launching web server", zap.Error(err))
		}
	},
}
