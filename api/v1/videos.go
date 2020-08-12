package v1

import (
	"context"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Videos
// @Summary Query videos
// @Description get latest videos
// @Accept  json
// @Produce  json
// @Success 200 {object} JSONResponse{data=[]ent.Video}
// @Failure 500 {object} JSONResponse
// @Router /v1/videos [get]
func getVideos(ctx *gin.Context) {
	ctx.JSON(200, JSONResponse{
		Success: true,
		Data:    []ent.Video{},
	})
}

// @Tags Videos
// @Summary Get a video
// @Description get one video
// @Accept  json
// @Produce  json
// @Param id path string true "Video ID" minlength(36) maxlength(36) validate(required)
// @Success 200 {object} JSONResponse{data=ent.Video}
// @Failure 404 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /v1/videos/{id} [get]
func getVideo(ctx *gin.Context) {
	ctx.JSON(200, JSONResponse{
		Success: true,
		Data:    new(ent.Video),
	})
}

// @Tags Videos
// @Summary Delete a video
// @Description Delete one video
// @Accept  json
// @Produce  json
// @Param id path string true "Video ID" minlength(36) maxlength(36) validate(required)
// @Success 200 {object} JSONResponse
// @Failure 404 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /v1/videos/{id} [delete]
func deleteVideo(ctx *gin.Context) {
	ctx.JSON(200, JSONResponse{
		Success: true,
	})
}

// @Tags Videos
// @Summary Create a video
// @Description Create a video
// @Accept  json
// @Produce  json
// @Success 200 {object} JSONResponse{data=ent.Video}
// @Failure 500 {object} JSONResponse
// @Router /v1/videos [post]
func createVideo(ctx *gin.Context) {
	v, err := db.Client.Video.
		Create().
		SetUUID(uuid.New().String()).
		SetTitle("test").
		SetStatus(schema.VideoStatusProcessing).
		Save(context.Background())
	if err != nil {
		ctx.JSON(500, JSONResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, JSONResponse{
		Success: true,
		Data:    v,
	})
}
