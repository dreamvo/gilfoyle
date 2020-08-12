package v1

import (
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/gin-gonic/gin"
)

// GetVideos godoc
// @Summary Query videos
// @Description get latest videos
// @Accept  json
// @Produce  json
// @Success 200 {object} JSONResult{data=[]ent.Video}
// @Failure 500 {object} JSONResult
// @Router /v1/videos [get]
func getVideos(ctx *gin.Context) {
	ctx.JSON(200, JSONResult{
		Data: []ent.Video{},
	})
}
