package api

import (
	"context"
	"errors"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/mediafile"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CreateMedia struct {
	Title string `json:"title" validate:"required,gte=1,lte=255" example:"Sheep Discovers How To Use A Trampoline"`
}

type UpdateMedia struct {
	CreateMedia
}

type MediasMetadata struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// @ID getAllMedias
// @Tags Medias
// @Summary Query medias
// @Description Get latest created medias
// @Produce  json
// @Success 200 {object} util.DataResponse{data=[]ent.Media,metadata=api.MediasMetadata}
// @Failure 500 {object} util.ErrorResponse
// @Router /medias [get]
// @Param limit query int false "Max number of results"
// @Param offset query int false "Number of results to ignore"
func (s *Server) getAllMedias(ctx *gin.Context) {
	limit := ctx.GetInt("limit")
	offset := ctx.GetInt("offset")

	total, err := s.db.Media.Query().Count(context.Background())
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	medias, err := s.db.Media.
		Query().
		Order(ent.Desc(media.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(context.Background())
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	util.NewData(ctx, http.StatusOK, medias, MediasMetadata{
		Total:  total,
		Offset: ctx.GetInt("offset"),
		Limit:  ctx.GetInt("limit"),
	})
}

// @ID getMedia
// @Tags Medias
// @Summary Get a media
// @Description Get one media
// @Produce  json
// @Param id path string true "Media ID" validate(required)
// @Success 200 {object} util.DataResponse{data=ent.Media}
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id} [get]
func (s *Server) getMedia(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, ErrInvalidUUID)
		return
	}

	v, err := s.db.Media.
		Query().
		Where(media.ID(parsedUUID)).
		WithMediaFiles().
		Only(context.Background())
	if v == nil {
		util.NewError(ctx, http.StatusNotFound, ErrResourceNotFound)
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	util.NewData(ctx, http.StatusOK, v, nil)
}

// @ID deleteMedia
// @Tags Medias
// @Summary Delete a media
// @Description Delete one media
// @Produce  json
// @Param id path string true "Media ID" validate(required)
// @Success 200 {object} util.DataResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id} [delete]
func (s *Server) deleteMedia(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := util.ValidateUUID(id)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, ErrInvalidUUID)
		return
	}

	v, _ := s.db.Media.Get(context.Background(), parsedUUID)
	if v == nil {
		util.NewError(ctx, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	if v.Status == media.StatusProcessing {
		util.NewError(ctx, http.StatusForbidden, errors.New("you can't delete a media while it's in processing state"))
		return
	}

	_, err = s.db.MediaFile.Delete().Where(mediafile.HasMediaWith(media.ID(parsedUUID))).Exec(context.Background())
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	err = s.storage.Delete(context.Background(), parsedUUID.String())
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	err = s.db.Media.DeleteOneID(parsedUUID).Exec(context.Background())
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	util.NewData(ctx, http.StatusOK, nil, nil)
}

// @ID createMedia
// @Tags Medias
// @Summary Create a media
// @Description Create a new media
// @Accept  json
// @Produce  json
// @Success 200 {object} util.DataResponse{data=ent.Media}
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias [post]
// @Param media body CreateMedia true "Media data" validate(required)
func (s *Server) createMedia(ctx *gin.Context) {
	var body CreateMedia

	if err := util.ValidateBody(ctx, &body); err != nil {
		util.NewValidationError(ctx, err)
		return
	}

	v, err := s.db.Media.
		Create().
		SetTitle(body.Title).
		SetStatus(schema.MediaStatusAwaitingUpload).
		Save(context.Background())
	if ent.IsValidationError(err) {
		util.NewError(ctx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	util.NewData(ctx, http.StatusOK, v, nil)
}

// @ID updateMedia
// @Tags Medias
// @Summary Update a media
// @Description Update an existing media
// @Accept  json
// @Produce  json
// @Success 200 {object} util.DataResponse{data=ent.Media}
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id} [patch]
// @Param id path string true "Media ID" validate(required)
// @Param media body UpdateMedia true "Media data" validate(required)
func (s *Server) updateMedia(ctx *gin.Context) {
	var body CreateMedia

	if err := util.ValidateBody(ctx, &body); err != nil {
		util.NewValidationError(ctx, err)
		return
	}

	id := ctx.Param("id")

	parsedUUID, err := util.ValidateUUID(id)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, ErrInvalidUUID)
		return
	}

	m, _ := s.db.Media.Get(context.Background(), parsedUUID)
	if m == nil {
		util.NewError(ctx, http.StatusNotFound, ErrResourceNotFound)
		return
	}

	m, err = s.db.Media.
		UpdateOneID(parsedUUID).
		SetTitle(body.Title).
		Save(context.Background())
	if ent.IsValidationError(err) {
		util.NewError(ctx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	util.NewData(ctx, http.StatusOK, m, nil)
}
