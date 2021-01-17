package gilfoyle

import (
	"context"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/ent/migrate"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func init() {
	// Register command
	rootCmd.AddCommand(serveCmd)

	// Register flags
	serveCmd.PersistentFlags().StringVar(&addr, "addr", "", "Interface binding for the web server")
	serveCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 3000, "HTTP port")
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Launch REST API",
	Example: "gilfoyle serve -p 3000 -c /app/config.yml",
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

		dbClient, err := db.NewClient(gilfoyle.Config.Services.DB)
		if err != nil {
			logger.Fatal("failed opening connection", zap.Error(err))
		}

		if !gilfoyle.Config.Settings.Debug {
			dbClient = dbClient.Debug()
		}

		// run the auto migration tool.
		if err := dbClient.Schema.Create(
			context.Background(),
			migrate.WithDropIndex(true),
			migrate.WithDropColumn(true),
			migrate.WithForeignKeys(true),
		); err != nil {
			logger.Fatal("failed creating schema resources", zap.Error(err))
		}

		logger.Info("Successfully executed database auto migration")

		s, err := gilfoyle.NewStorage(config.StorageDriver(gilfoyle.Config.Storage.Driver))
		if err != nil {
			logger.Fatal("Error initializing storage backend", zap.Error(err))
		}

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
		defer func() { _ = w.Close() }()

		err = w.Init()
		if err != nil {
			logger.Fatal("Failed to initialize worker queues", zap.Error(err))
		}

		server := api.NewServer(api.Options{
			Logger:   logger,
			Worker:   w,
			Database: dbClient,
			Config:   gilfoyle.Config,
			Storage:  s,
		})

		// Launch web server
		if err := server.Listen(fmt.Sprintf("%s:%d", addr, httpPort)); err != nil {
			logger.Fatal("error while launching web server", zap.Error(err))
		}
	},
}
