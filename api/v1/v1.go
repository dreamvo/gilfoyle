package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	defaultItemsPerPage = 50
	maxItemsPerPage     = 100
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
			videos.PATCH(":id", updateVideo)
			videos.POST(":id/upload", uploadVideoFile)
		}
	}

	return v1
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
