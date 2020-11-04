//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go
package api

import (
	_ "github.com/dreamvo/gilfoyle/api/docs"
	"github.com/dreamvo/gilfoyle/config"
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
// @description Cloud-native media hosting & streaming server for businesses.
// @version v1
// @host demo-v1.gilfoyle.dreamvo.com
// @BasePath /
// @schemes http https
// @license.name GNU General Public License v3.0
// @license.url https://github.com/dreamvo/gilfoyle/blob/master/LICENSE

// RegisterRoutes adds routes to a given router instance
func RegisterRoutes(r *gin.Engine) *gin.Engine {
	r.GET("/health", healthCheckHandler)

	medias := r.Group("/medias")
	{
		medias.GET("", paginateHandler, getAllMedias)
		medias.GET(":id", getMedia)
		medias.DELETE(":id", deleteMedia)
		medias.POST("", createMedia)
		medias.PATCH(":id", updateMedia)
		medias.POST(":id/upload", uploadMediaFile)
	}

	if config.GetConfig().Settings.ExposeSwaggerUI {
		// Register swagger docs handler
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

// @ID checkHealth
// @Tags health
// @Summary Check service status
// @Description Check for the health of the service
// @Produce  json
// @Success 200 {object} interface{}
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
