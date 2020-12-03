package restapi

import (
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckResponse struct {
	Tag    string `json:"tag"`
	Commit string `json:"commit"`
}

func RegisterStaticRoutes(r *gin.Engine) *gin.Engine {
	// TODO(sundowndev): fix static assets serving
	r.Static("/", "../../dashboard/ui/dist/")

	return r
}

func RegisterAPIRoutes(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	{
		api.GET("/healthz", healthCheckHandler)
	}

	return r
}

func healthCheckHandler(ctx *gin.Context) {
	httputils.NewResponse(ctx, http.StatusOK, HealthCheckResponse{
		Tag:    config.Version,
		Commit: config.Commit,
	})
}
