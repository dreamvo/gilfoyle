package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) *gin.RouterGroup {
	v1 := r.Group("/v1")
	{
		videos := v1.Group("/videos")
		{
			videos.GET("", getVideos)
			videos.GET(":id", getVideo)
			videos.DELETE(":id", deleteVideo)
			videos.POST("", createVideo)
		}
	}

	return v1
}
