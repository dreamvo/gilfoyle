package api

import (
	"github.com/gin-gonic/gin"
)

// @ID getMetrics
// @Tags Metrics
// @Summary Get instance metrics
// @Description Get metrics about this Gilfoyle instance
// @Produce  json
// @Success 200 {object} util.DataResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /metrics [get]
func getMetrics(ctx *gin.Context) {
	ctx.Status(200)
}
