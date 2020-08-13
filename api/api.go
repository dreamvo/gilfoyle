//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go
package api

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/api/docs"
	"github.com/dreamvo/gilfoyle/api/v1"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @license.name GNU General Public License v3.0
// @license.url https://github.com/dreamvo/gilfoyle/blob/master/LICENSE

// RegisterRoutes runs a REST API web server
func RegisterRoutes(r *gin.Engine, port int) *gin.Engine {
	docs.SwaggerInfo.Title = "Gilfoyle server"
	docs.SwaggerInfo.Description = " Video streaming server backed by decentralized filesystem."
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/health", healthcheckHandler)

	v1.RegisterRoutes(r)

	// register swagger docs handler
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// @Tags health
// @Summary Check service status
// @Success 200
// @Router /health [get]
func healthcheckHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(200)
}
