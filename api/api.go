//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go
package api

import (
	"github.com/dreamvo/gilfoyle/api/v1"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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
func RegisterRoutes(r *gin.Engine, port int, serveDocs bool) *gin.Engine {
	r.GET("/health", healthcheckHandler)

	v1.RegisterRoutes(r)

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
func healthcheckHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(200)
}
