package v1

import (
	"context"
	"fmt"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/ent/video"
	"github.com/dreamvo/gilfoyle/httputils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type createVideoBody struct {
	Title string `json:"title"`
}

// @Tags videos
// @Summary Query videos
// @Description get latest videos
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=[]ent.Video}
// @Failure 500 {object} httputils.ErrorResponse
// @Router /v1/videos [get]
// @Param limit query int false "Max number of results" minimum(1) maximum(100)
// @Param offset query int false "Number of results to ignore" minimum(0)
func getVideos(ctx *gin.Context) {
	limit := ctx.GetInt("limit")
	offset := ctx.GetInt("offset")

	videos, err := db.Client.Video.
		Query().
		Order(ent.Desc(video.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(context.Background())
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	httputils.NewData(ctx, http.StatusOK, videos)
}

// @Tags videos
// @Summary Get a video
// @Description get one video
// @Produce  json
// @Param id path string true "Video ID" minlength(36) maxlength(36) validate(required)
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /v1/videos/{id} [get]
func getVideo(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, fmt.Errorf("invalid UUID provided"))
		return
	}

	v, err := db.Client.Video.Get(context.Background(), parsedUUID)
	if err != nil {
		httputils.NewError(ctx, http.StatusNotFound, err)
		return
	}

	httputils.NewData(ctx, http.StatusOK, v)
}

// @Tags videos
// @Summary Delete a video
// @Description Delete one video
// @Produce  json
// @Param id path string true "Video ID" minlength(36) maxlength(36) validate(required)
// @Success 200 {object} httputils.DataResponse
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /v1/videos/{id} [delete]
func deleteVideo(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, fmt.Errorf("invalid UUID provided"))
		return
	}

	err = db.Client.Video.DeleteOneID(parsedUUID).Exec(context.Background())
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	httputils.NewData(ctx, http.StatusOK, nil)
}

// @Tags videos
// @Summary Create a video
// @Description Create a new video
// @Accept  json
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /v1/videos [post]
// @Param title body string true "Video title" minlength(1) maxlength(255) validate(required)
func createVideo(ctx *gin.Context) {
	var body createVideoBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	v, err := db.Client.Video.
		Create().
		SetTitle(body.Title).
		SetStatus(schema.VideoStatusProcessing).
		Save(context.Background())
	if ent.IsValidationError(err) {
		httputils.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	httputils.NewData(ctx, http.StatusOK, v)
}

// @Tags videos
// @Summary Upload a video file
// @Description Upload a new video file for a given video ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /v1/videos/{id}/upload [post]
// @Param id path string true "Video ID" minlength(36) maxlength(36) validate(required)
// @Param file formData file true "Video file"
func uploadVideoFile(ctx *gin.Context) {
	ctx.Status(200)
}
