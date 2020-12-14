//go:generate go run github.com/rakyll/statik -src ./ui/dist -include=*.png,*.ico,*.html,*.css,*.js
package dashboard

import (
	_ "github.com/dreamvo/gilfoyle/dashboard/statik"
	"github.com/dreamvo/gilfoyle/logging"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Server struct {
	router   *gin.Engine
	logger   logging.ILogger
	endpoint string
}

type HealthCheckResponse struct {
	Tag    string `json:"tag"`
	Commit string `json:"commit"`
}

func NewServer(logger logging.ILogger, endpoint string) *Server {
	router := gin.New()
	s := &Server{
		router:   router,
		logger:   logger,
		endpoint: endpoint,
	}

	s.router.Use(func(ctx *gin.Context) {
		ctx.Next()

		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		errorMsg := ctx.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logMsg := logger.With(
			zap.String("Method", ctx.Request.Method),
			zap.String("Path", path),
			zap.Int("StatusCode", ctx.Writer.Status()),
			zap.Int("BodySize", ctx.Writer.Size()),
			zap.String("UserAgent", ctx.Request.UserAgent()),
		)

		if errorMsg != "" {
			logMsg.Error("Incoming HTTP Request",
				zap.String("ErrorMessage", errorMsg),
			)
			return
		}

		logMsg.Info("Incoming HTTP Request")
	})

	s.router.Use(gin.Recovery())

	registerAPIRoutes(s)
	registerStaticRoutes(s)

	return s
}

func registerStaticRoutes(s *Server) *Server {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// TODO(sundowndev): fix static assets serving
	s.router.StaticFS("/app/", statikFS)

	s.router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/app")
	})

	return s
}

func registerAPIRoutes(s *Server) *Server {
	api := s.router.Group("/api")
	{
		api.GET("/proxy/*path", s.proxyHandler)
	}

	return s
}

func (s *Server) proxyHandler(ctx *gin.Context) {
	//path := fmt.Sprintf("%s/%s", s.endpoint, ctx.Param("path"))
	ctx.Status(200)
}

func (s *Server) Listen(addr ...string) error {
	return s.router.Run(addr...)
}
