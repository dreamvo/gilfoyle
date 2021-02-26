package api

import (
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckResponse struct {
	Tag             string `json:"tag"`
	Commit          string `json:"commit"`
	Debug           bool   `json:"debug"`
	StorageDriver   string `json:"storage_driver"`
	MaxFileSize     int64  `json:"max_file_size"`
	DatabaseDialect string `json:"database_dialect"`
}

// @ID checkHealth
// @Tags Instance
// @Summary Check service status
// @Description Check for the health of the service
// @Produce  json
// @Success 200 {object} util.DataResponse{data=HealthCheckResponse}
// @Router /healthz [get]
func (s *Server) healthCheckHandler(ctx *gin.Context) {
	util.NewResponse(ctx, http.StatusOK, HealthCheckResponse{
		Tag:             config.Version,
		Commit:          config.Commit,
		Debug:           s.config.Settings.Debug,
		StorageDriver:   s.config.Storage.Driver,
		MaxFileSize:     s.config.Settings.MaxFileSize,
		DatabaseDialect: s.config.Services.DB.Dialect,
	})
}

// @ID getPrometheusMetrics
// @Tags Instance
// @Summary Get instance metrics
// @Description Get Prometheus exported metrics about this Gilfoyle instance
// @Produce  text/plain
// @Success 200 {string} string
// @Router /metricsz [get]
