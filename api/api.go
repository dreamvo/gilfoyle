//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go
package api

import (
	_ "github.com/dreamvo/gilfoyle/api/docs"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"strconv"
)

const (
	defaultItemsPerPage = 50
	maxItemsPerPage     = 100
	ErrInvalidUUID      = "invalid UUID provided"
	ErrResourceNotFound = "resource not found"
)

// @title Gilfoyle server
// @description Video streaming server backed by decentralized filesystem.
// @version 0.1-beta
// @host localhost:8080
// @BasePath /
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/dreamvo/gilfoyle/blob/master/LICENSE

// RegisterRoutes adds routes to a given router instance
func RegisterRoutes(r *gin.Engine, serveDocs bool) *gin.Engine {
	r.GET("/health", healthCheckHandler)

	videos := r.Group("/videos")
	{
		videos.GET("", paginateHandler, getVideos)
		videos.GET(":id", getVideo)
		videos.DELETE(":id", deleteVideo)
		videos.POST("", createVideo)
		videos.PATCH(":id", updateVideo)
		videos.POST(":id/upload", uploadVideoFile)
	}

	if serveDocs {
		// register swagger docs handler
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

// @Tags health
// @Summary Check service status
// @Success 200
// @Router /health [get]
func healthCheckHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(200)
}

func paginateHandler(ctx *gin.Context) {
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
