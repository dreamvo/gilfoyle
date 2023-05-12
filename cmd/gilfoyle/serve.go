package gilfoyle

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent/migrate"
	"github.com/dreamvo/gilfoyle/logging"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
		logger, err := logging.NewLogger(cfg.Settings.Debug, true)
		if err != nil {
			log.Fatal(err)
		}

		logger.Info("Initializing API server")
		logger.Info("Environment", zap.Bool("debug", cfg.Settings.Debug))

		if !cfg.Settings.Debug {
			gin.SetMode(gin.ReleaseMode)
		} else {
			_ = os.Setenv("PGSSLMODE", "disable")
			gin.SetMode(gin.DebugMode)
		}

		dbClient, err := db.NewClient(cfg.Services.DB)
		if err != nil {
			logger.Fatal("failed opening connection", zap.Error(err))
		}

		if !cfg.Settings.Debug {
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

		s, err := gilfoyle.NewStorage(*cfg)
		if err != nil {
			logger.Fatal("Error initializing storage backend", zap.Error(err))
		}

		w, err := worker.New(worker.Options{
			Host:        cfg.Services.RabbitMQ.Host,
			Port:        cfg.Services.RabbitMQ.Port,
			Username:    cfg.Services.RabbitMQ.Username,
			Password:    cfg.Services.RabbitMQ.Password,
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
			Config:   *cfg,
			Storage:  s,
		})

		// Launch web server
		if err := server.Listen(fmt.Sprintf("%s:%d", addr, httpPort)); err != nil {
			logger.Fatal("error while launching web server", zap.Error(err))
		}
	},
}
