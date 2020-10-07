package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/ent/video"
	"github.com/dreamvo/gilfoyle/httputils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
)


type CreateVideo struct {
	Title string `json:"title" validate:"required,gte=1,lte=255" example:"Sheep Discovers How To Use A Trampoline"`
}

type UpdateVideo struct {
	CreateVideo
}

// @ID getAllVideos
// @Tags videos
// @Summary Query videos
// @Description get latest videos
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=[]ent.Video}
// @Failure 500 {object} httputils.ErrorResponse
// @Router /videos [get]
// @Param limit query int false "Max number of results"
// @Param offset query int false "Number of results to ignore"
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
		httputils.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	httputils.NewData(ctx, http.StatusOK, videos)
}

// @ID getVideo
// @Tags videos
// @Summary Get a video
// @Description get one video
// @Produce  json
// @Param id path string true "Video ID" validate(required)
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /videos/{id} [get]
func getVideo(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, fmt.Errorf(ErrInvalidUUID))
		return
	}

	v, err := db.Client.Video.Get(context.Background(), parsedUUID)
	if v == nil {
		httputils.NewError(ctx, http.StatusNotFound, errors.New(ErrResourceNotFound))
		return
	}
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	httputils.NewData(ctx, http.StatusOK, v)
}

// @ID deleteVideo
// @Tags videos
// @Summary Delete a video
// @Description Delete one video
// @Produce  json
// @Param id path string true "Video ID" validate(required)
// @Success 200 {object} httputils.DataResponse
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /videos/{id} [delete]
func deleteVideo(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, fmt.Errorf(ErrInvalidUUID))
		return
	}

	v, _ := db.Client.Video.Get(context.Background(), parsedUUID)
	if v == nil {
		httputils.NewError(ctx, http.StatusNotFound, errors.New(ErrResourceNotFound))
		return
	}

	err = db.Client.Video.DeleteOneID(parsedUUID).Exec(context.Background())
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	httputils.NewData(ctx, http.StatusOK, nil)
}

// @ID createVideo
// @Tags videos
// @Summary Create a video
// @Description Create a new video
// @Accept  json
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /videos [post]
// @Param video body CreateVideo true "Video data" validate(required)
func createVideo(ctx *gin.Context) {
	err := validator.New().StructCtx(ctx, CreateVideo{})
	if err != nil {
		httputils.NewValidationError(ctx, http.StatusBadRequest, err)
		return
	}

	var body CreateVideo
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
		httputils.NewError(ctx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	httputils.NewData(ctx, http.StatusOK, v)
}

// @ID updateVideo
// @Tags videos
// @Summary Update a video
// @Description Update an existing video
// @Accept  json
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /videos/{id} [patch]
// @Param id path string true "Video ID" validate(required)
// @Param video body UpdateVideo true "Video data" validate(required)
func updateVideo(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, fmt.Errorf(ErrInvalidUUID))
		return
	}

	v, _ := db.Client.Video.Get(context.Background(), parsedUUID)
	if v == nil {
		httputils.NewError(ctx, http.StatusNotFound, errors.New(ErrResourceNotFound))
		return
	}

	var body UpdateVideo
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httputils.NewError(ctx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}

	v, err = db.Client.Video.
		UpdateOneID(parsedUUID).
		SetTitle(body.Title).
		Save(context.Background())
	if ent.IsValidationError(err) {
		httputils.NewError(ctx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}
	if err != nil {
		httputils.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	httputils.NewData(ctx, http.StatusOK, v)
}

// @ID uploadVideoFile
// @Tags videos
// @Summary Upload a video file
// @Description Upload a new video file for a given video ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} httputils.DataResponse{data=ent.Video}
// @Failure 404 {object} httputils.ErrorResponse
// @Failure 400 {object} httputils.ErrorResponse
// @Failure 500 {object} httputils.ErrorResponse
// @Router /videos/{id}/upload [post]
// @Param id path string true "Video ID" validate(required)
// @Param file formData file true "Video file"
func uploadVideoFile(ctx *gin.Context) {
	ctx.Status(200)
}
