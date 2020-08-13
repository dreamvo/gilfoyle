package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	v1 := r.Group("/v1")
	{
		videos := v1.Group("/videos")
		{
			videos.GET("", paginateHandler, getVideos)
			videos.GET(":id", getVideo)
			videos.DELETE(":id", deleteVideo)
			videos.POST("", createVideo)
		}
	}

	return v1
}

func paginateHandler(ctx *gin.Context) {
	limit := ctx.Query("limit")
	limitInt, err := strconv.ParseInt(limit, 10, 64)

	if err != nil || limitInt > 100 {
		limitInt = 50
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
