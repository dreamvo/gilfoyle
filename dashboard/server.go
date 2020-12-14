//go:generate go run github.com/rakyll/statik -src ./ui/dist -include=*.png,*.ico,*.html,*.css,*.js
package dashboard

import (
	"fmt"
	_ "github.com/dreamvo/gilfoyle/dashboard/statik"
	"github.com/dreamvo/gilfoyle/logging"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	staticRootPath string = "/app"
)

type Server struct {
	router   *gin.Engine
	logger   logging.ILogger
	endpoint string
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
		s.logger.Fatal("register static routes", zap.Error(err))
	}

	s.router.StaticFS(staticRootPath, statikFS)

	s.router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, staticRootPath)
	})

	return s
}

func registerAPIRoutes(s *Server) *Server {
	api := s.router.Group("/api")
	{
		api.Any("/proxy/*path", s.proxyHandler)
	}

	return s
}

func (s *Server) proxyHandler(ctx *gin.Context) {
	pathSegments := strings.Split(ctx.Param("path"), "/")
	path := strings.Join(pathSegments[1:], "/")
	fullPath := fmt.Sprintf("%s/%s", s.endpoint, path)

	req, err := http.NewRequestWithContext(ctx, ctx.Request.Method, fullPath, ctx.Request.Body)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
	}

	ctx.DataFromReader(res.StatusCode, res.ContentLength, res.Header.Get("Content-Type"), res.Body, map[string]string{})
}

func (s *Server) Listen(addr ...string) error {
	return s.router.Run(addr...)
}
