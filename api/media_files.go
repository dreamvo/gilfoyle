package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent"
	_ "github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/vansante/go-ffprobe.v2"
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

// @ID uploadVideo
// @Tags Medias
// @Summary Upload a video file
// @Description Upload a new video file for a given media ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} util.DataResponse{data=FileFormat}
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

	if file.Size > s.config.Settings.MaxFileSize {
		util.NewError(ctx, http.StatusBadRequest, fmt.Errorf("uploaded file's size exceed limit of %v", s.config.Settings.MaxFileSize))
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

	tx, err := s.db.Tx(context.Background())
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	_, err = tx.Media.
		UpdateOneID(m.ID).
		SetStatus(media.StatusProcessing).
		SetOriginalFilename(transcoding.OriginalFileName).
		Save(context.Background())
	if ent.IsValidationError(err) {
		rollbackWithError(ctx, tx, http.StatusBadRequest, errors.Unwrap(err))
		return
	}
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, errors.Unwrap(err))
		return
	}

	videoStream := data.StreamType(ffprobe.StreamVideo)
	if len(videoStream) != 1 {
		util.NewError(ctx, http.StatusBadRequest, errors.New("uploaded media must have 1 video stream"))
		return
	}

	err = tx.Commit()
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ch, err := s.worker.Client.Channel()
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	format := strings.Split(data.Format.FormatName, ",")[0]
	fps := int(transcoding.ParseFrameRates(data.Streams[0].RFrameRate))

	RenditionsCount := 0

	for _, r := range s.config.Settings.Encoding.Renditions {
		// Ignore resolutions higher than original
		if r.Width > data.FirstVideoStream().Width && r.Height > data.FirstVideoStream().Height {
			continue
		}

		if r.Framerate != 0 {
			fps = r.Framerate
		}

		err = worker.VideoTranscodingProducer(ch, worker.VideoTranscodingParams{
			MediaUUID: m.ID,
			OriginalFile: transcoding.OriginalFile{
				Format:          format,
				FrameRate:       uint8(fps),
				DurationSeconds: data.Format.DurationSeconds,
				Filepath:        path,
			},
			RenditionName:      r.Name,
			FrameRate:          fps,
			VideoWidth:         r.Width,
			VideoHeight:        r.Height,
			AudioCodec:         r.AudioCodec,
			VideoCodec:         r.VideoCodec,
			Crf:                20,
			KeyframeInterval:   48,
			HlsSegmentDuration: 4,
			HlsPlaylistType:    "vod",
			VideoBitRate:       r.VideoBitrate,
			AudioBitrate:       r.AudioBitrate,
			TargetBandwidth:    uint64(r.VideoBitrate + r.AudioBitrate),
		})
		if err != nil {
			util.NewError(ctx, http.StatusInternalServerError, err)
			return
		}

		RenditionsCount++
	}

	err = worker.MediaProcessingCallbackProducer(ch, worker.MediaProcessingCallbackParams{
		MediaUUID:       m.ID,
		MediaFilesCount: RenditionsCount,
	})
	if err != nil {
		util.NewError(ctx, http.StatusInternalServerError, err)
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
	}, nil)
}

// @Deprecated
// @ID uploadAudio
// @Tags Medias
// @Summary Upload a audio file
// @Description Upload a new audio file for a given media ID
// @Accept  multipart/form-data
// @Produce  json
// @Success 200 {object} util.DataResponse{data=FileFormat}
// @Failure 404 {object} util.ErrorResponse
// @Failure 400 {object} util.ErrorResponse
// @Failure 500 {object} util.ErrorResponse
// @Router /medias/{id}/upload/audio [post]
// @Param id path string true "Media identifier" validate(required)
// @Param file formData file true "Audio file"
func (s *Server) uploadAudioFile(ctx *gin.Context) {
	ctx.Status(200)
}
