package api

import (
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckResponse struct {
	Tag    string `json:"tag"`
	Commit string `json:"commit"`
}

// @ID checkHealth
// @Tags Instance
// @Summary Check service status
// @Description Check for the health of the service
// @Produce  json
// @Success 200 {object} util.DataResponse{data=HealthCheckResponse}
// @Router /healthz [get]
func healthCheckHandler(ctx *gin.Context) {
	util.NewResponse(ctx, http.StatusOK, HealthCheckResponse{
		Tag:    config.Version,
		Commit: config.Commit,
	})
}

// @ID getMetrics
// @Tags Instance
// @Summary Get instance metrics
// @Description Get metrics about this Gilfoyle instance
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /metricsz [get]
func getMetrics(ctx *gin.Context) {
	ctx.Status(200)
}
