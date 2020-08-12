package v1

import (
	"github.com/gin-gonic/gin"
)

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Router(r *gin.Engine) *gin.RouterGroup {
	v1 := r.Group("/v1")
	{
		videos := v1.Group("/videos")
		{
			videos.GET("", getVideos)
		}
	}

	return v1
}
