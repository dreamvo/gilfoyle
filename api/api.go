//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go --parseDependency
package api

import (
	"errors"
	_ "github.com/dreamvo/gilfoyle/api/docs"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/logging"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	defaultItemsPerPage = 50
	maxItemsPerPage     = 100
)

var (
	ErrInvalidUUID      = errors.New("invalid UUID provided")
	ErrResourceNotFound = errors.New("resource not found")
)

type Options struct {
	Database *ent.Client
	Config   config.Config
	Storage  storage.Storage
	Worker   *worker.Worker
	Logger   logging.ILogger
}

type Server struct {
	router  *gin.Engine
	db      *ent.Client
	worker  *worker.Worker
	logger  logging.ILogger
	storage storage.Storage
	config  config.Config
}

// @title Gilfoyle server
// @description Cloud-native media hosting & streaming server for businesses.
// @version v1
// @host demo-v1.gilfoyle.dreamvo.com
// @BasePath /
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/dreamvo/gilfoyle/blob/master/LICENSE

func NewServer(opts Options) *Server {
	s := &Server{
		router:  gin.New(),
		storage: opts.Storage,
		db:      opts.Database,
		config:  opts.Config,
		logger:  opts.Logger,
		worker:  opts.Worker,
	}
	registerMiddlewares(s)
	registerRoutes(s)
	return s
}

// registerMiddlewares adds middlewares to a given router instance
func registerMiddlewares(s *Server) {
	s.router.Use(gin.Recovery())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	// TODO(sundowndev): update Gin to enable this feature. See https://github.com/gin-gonic/gin/commits/master/recovery.go
	//r.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
	//	if err, ok := recovered.(string); ok {
	//		util.NewError(ctx, http.StatusInternalServerError, errors.New(err))
	//	}
	//	util.NewError(ctx, http.StatusInternalServerError, errors.New("an unexpected error occurred"))
	//}))

	s.router.Use(func(ctx *gin.Context) {
		ctx.Next()

		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		errorMsg := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log := s.logger.With(
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
}

// registerRoutes adds routes to a given router instance
func registerRoutes(s *Server) {
	s.router.GET("/healthz", s.healthCheckHandler)

	p := ginprometheus.NewPrometheus("gin")
	p.MetricsPath = "/metricsz"
	p.Use(s.router)

	medias := s.router.Group("/medias")
	{
		medias.GET("", s.paginateHandler, s.getAllMedias)
		medias.GET(":id", s.getMedia)
		medias.DELETE(":id", s.deleteMedia)
		medias.POST("", s.createMedia)
		medias.PATCH(":id", s.updateMedia)
		medias.POST(":id/upload/video", s.uploadVideoFile)
		medias.POST(":id/upload/audio", s.uploadAudioFile)
		medias.GET(":id/attachments", s.getMediaAttachments)
		medias.POST(":id/attachments", s.addMediaAttachments)
		medias.DELETE(":id/attachments/:attachment_id", s.deleteMediaAttachments)
		medias.GET(":id/stream/:preset", s.streamMedia)
	}

	if s.config.Settings.ExposeSwaggerUI {
		// Register swagger docs handler
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	s.router.Use(func(ctx *gin.Context) {
		util.NewError(ctx, http.StatusNotFound, errors.New("resource not found"))
	})
}

func (s *Server) Listen(addr ...string) error {
	defer func() { _ = s.db.Close() }()
	return s.router.Run(addr...)
}

func (s *Server) paginateHandler(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt, err := strconv.ParseInt(limit, 10, 64)

	if err != nil || limitInt > maxItemsPerPage {
		limitInt = defaultItemsPerPage
	}

	offset := ctx.Query("offset")
	offsetInt, err := strconv.ParseInt(offset, 10, 64)

	if err != nil {
		offsetInt = 0
	}

	ctx.Set("limit", int(limitInt))
	ctx.Set("offset", int(offsetInt))
	ctx.Next()
}

func rollbackWithError(ctx *gin.Context, tx *ent.Tx, statusCode int, err error) {
	if txErr := tx.Rollback(); txErr != nil {
		util.NewError(ctx, statusCode, txErr)
	}

	util.NewError(ctx, statusCode, err)
}
