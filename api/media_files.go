package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/api/util"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/vansante/go-ffprobe.v2"
	"net/http"
	"os"
	"time"
)

type FileFormat struct {
	Filename         string  `json:"filename"`
	NBStreams        int     `json:"nb_streams"`
	NBPrograms       int     `json:"nb_programs"`
	FormatName       string  `json:"format_name"`
	FormatLongName   string  `json:"format_long_name"`
	StartTimeSeconds float64 `json:"start_time,string"`
	DurationSeconds  float64 `json:"duration,string"`
	Size             string  `json:"size"`
	BitRate          string  `json:"bit_rate"`
	ProbeScore       int     `json:"probe_score"`
}

// @ID uploadMediaFile
// @Tags Medias
// @Summary Upload a media file
// @Description Upload a new media file for a given media ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} util.DataResponse{data=FileFormat}
// @Failure 404 {object} util.ErrorResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id}/upload [post]
// @Param id path string true "Media identifier" validate(required)
// @Param file formData file true "Media file"
func uploadMediaFile(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, fmt.Errorf(ErrInvalidUUID))
		return
	}

	v, err := db.Client.Media.Get(context.Background(), parsedUUID)
	if v == nil {
		util.NewError(ctx, http.StatusNotFound, errors.New("media could not be found"))
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	path := fmt.Sprintf("%s/%s", id, "original")
	stat, _ := gilfoyle.Storage.Stat(context.Background(), path)
	if stat != nil {
		util.NewError(ctx, http.StatusBadRequest, errors.New("a file already exists for this media"))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if file.Size > gilfoyle.Config.Settings.MaxFileSize {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("uploaded file's size exceed limit of %v", gilfoyle.Config.Settings.MaxFileSize))
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error opening uploaded file: %e", err))
		return
	}

	ctxWithTimeout, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeReader(ctxWithTimeout, fileReader)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error analyzing uploaded file: %e", err))
		return
	}

	tmpPath := fmt.Sprintf("%s/%s", os.TempDir(), uuid.New().String())

	if err = ctx.SaveUploadedFile(file, tmpPath); err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error saving temporary file: %s", err))
		return
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error opening temporary file: %s", err))
		return
	}

	if err = gilfoyle.Storage.Save(ctx, f, path); err != nil {
		util.NewError(ctx, http.StatusInternalServerError, fmt.Errorf("error saving uploaded file: %s", err))
		return
	}

	util.NewData(ctx, http.StatusOK, FileFormat{
		Filename:         data.Format.Filename,
		NBStreams:        data.Format.NBStreams,
		NBPrograms:       data.Format.NBPrograms,
		FormatName:       data.Format.FormatName,
		FormatLongName:   data.Format.FormatLongName,
		StartTimeSeconds: data.Format.StartTimeSeconds,
		DurationSeconds:  data.Format.DurationSeconds,
		Size:             data.Format.Size,
		BitRate:          data.Format.BitRate,
		ProbeScore:       data.Format.ProbeScore,
	})
}
