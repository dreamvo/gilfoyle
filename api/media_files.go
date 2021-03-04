package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
)

type UploadMediaFileResponse struct {
	Success bool `json:"success"`
}

// @ID uploadVideo
// @Tags Medias
// @Summary Upload a video file
// @Description Upload a new video file for a given media ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} util.DataResponse{data=UploadMediaFileResponse}
// @Failure 404 {object} util.ErrorResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id}/upload/video [post]
// @Param id path string true "Media identifier" validate(required)
// @Param file formData file true "Video file"
func (s *Server) uploadVideoFile(ctx *gin.Context) {
	parsedUUID, err := util.ValidateUUID(ctx.Param("id"))
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, ErrInvalidUUID)
		return
	}

	m, err := s.db.Media.Get(context.Background(), parsedUUID)
	if m == nil {
		util.NewError(ctx, http.StatusNotFound, errors.New("media could not be found"))
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	if m.Status != media.StatusAwaitingUpload {
		util.NewError(ctx, http.StatusBadRequest, errors.New("a file already exists for this media"))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if file.Size > gilfoyle.Config.Settings.MaxFileSize {
		util.NewError(ctx, http.StatusBadRequest, fmt.Errorf("uploaded file's size exceed limit of %v", gilfoyle.Config.Settings.MaxFileSize))
		return
	}

	tmpPath := filepath.Join(os.TempDir(), uuid.New().String())

	if err = ctx.SaveUploadedFile(file, tmpPath); err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error saving temporary file: %s", err))
		return
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error opening temporary file: %s", err))
		return
	}

	path := fmt.Sprintf("%s/%s", parsedUUID.String(), transcoding.OriginalFileName)

	if err = s.storage.Save(ctx, f, path); err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error saving uploaded file: %s", err))
		return
	}

	_, err = s.db.Media.
		UpdateOneID(m.ID).
		SetStatus(media.StatusProcessing).
		SetOriginalFilename(transcoding.OriginalFileName).
		Save(context.Background())
	if ent.IsValidationError(err) {
		util.NewError(ctx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	ch, err := s.worker.Client.Channel()
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	err = worker.EncodingEntrypointProducer(ch, worker.EncodingEntrypointParams{
		MediaUUID: m.ID,
	})
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	util.NewData(ctx, http.StatusOK, UploadMediaFileResponse{Success: true}, nil)
}

// @Deprecated
// @ID uploadAudio
// @Tags Medias
// @Summary Upload a audio file
// @Description Upload a new audio file for a given media ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} util.DataResponse{data=UploadMediaFileResponse}
// @Failure 404 {object} util.ErrorResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id}/upload/audio [post]
// @Param id path string true "Media identifier" validate(required)
// @Param file formData file true "Audio file"
func (s *Server) uploadAudioFile(ctx *gin.Context) {
	ctx.Status(200)
}
