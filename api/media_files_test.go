package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/mediafile"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/dreamvo/gilfoyle/x/testutils"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func removeDir(path string) {
	_ = os.RemoveAll(path)
}

func TestMediaFiles(t *testing.T) {
	var r *gin.Engine

	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() { _ = dbClient.Close() }()

	cfg, err := gilfoyle.NewConfig()
	if err != nil {
		t.Error(err)
	}

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer removeDir(gilfoyle.Config.Storage.Filesystem.DataPath)

	storageDriver, err := gilfoyle.NewStorage(storage.Filesystem)
	if err != nil {
		t.Error(err)
	}

	container := testutils.CreateRabbitMQContainer(t, "guest", "guest")
	defer testutils.StopContainer(t, container)

	w, err := worker.New(worker.Options{
		Host:        container.Host,
		Port:        container.DefaultPort(),
		Username:    "guest",
		Password:    "guest",
		Logger:      zap.NewExample(),
		Concurrency: 1,
	})
	if err != nil {
		t.Error(err)
	}
	defer testutils.CloseWorker(t, w)

	err = w.Init()
	if err != nil {
		t.Error(err)
	}

	s := NewServer(Options{
		Database: dbClient,
		Config:   *cfg,
		Storage:  storageDriver,
		Worker:   w,
		Logger:   zap.NewExample(),
	})
	r = s.router

	t.Run("POST /medias/:id/upload/video", func(t *testing.T) {
		t.Run("should upload file and return probe", func(t *testing.T) {
			m, _ := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)

			filePath := "../x/testutils/fixtures/SampleVideo_1280x720_1mb.mp4"

			file, err := os.Open(filePath)
			assert.NoError(t, err)
			defer file.Close()

			part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
			assert.NoError(t, err)

			_, err = io.Copy(part1, file)
			assert.NoError(t, err)

			err = writer.Close()
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/medias/%s/upload/video", m.ID), payload)
			assert.NoError(t, err)

			req.Header.Add("Content-Type", "multipart/form-data")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			var body util.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			stat, err := os.Stat(filepath.Join(
				gilfoyle.Config.Storage.Filesystem.DataPath,
				m.ID.String(),
				transcoding.OriginalFileName,
			))
			assert.NoError(t, err)
			assert.Equal(t, int64(1055736), stat.Size())

			assert.Equal(t, 200, res.Result().StatusCode)
			assert.Equal(t, map[string]interface{}{
				"bit_rate":         "",
				"duration":         "5.312",
				"filename":         "pipe:",
				"format_long_name": "QuickTime / MOV",
				"format_name":      "mov,mp4,m4a,3gp,3g2,mj2",
				"nb_programs":      float64(0),
				"nb_streams":       float64(2),
				"probe_score":      float64(100),
				"size":             "",
				"start_time":       "0",
			}, body.Data)

			m, _ = dbClient.Media.Get(context.Background(), m.ID)

			assert.Equal(t, media.StatusProcessing, m.Status)

			mediaFile, _ := dbClient.MediaFile.
				Query().
				Where(mediafile.MediaTypeEQ(schema.MediaFileTypeVideo)).
				Only(context.Background())

			assert.Equal(t, int8(25), mediaFile.Framerate)
			assert.Equal(t, 5.312, mediaFile.DurationSeconds)
			assert.Equal(t, int16(1280), mediaFile.ScaledWidth)
			assert.Equal(t, mediafile.EncoderPreset(schema.MediaFileEncoderPresetOriginal), mediaFile.EncoderPreset)
			assert.Equal(t, int64(1205959), mediaFile.VideoBitrate)
			assert.Equal(t, mediafile.MediaType(schema.MediaFileTypeVideo), mediaFile.MediaType)

			ch, err := w.Client.Channel()
			assert.NoError(t, err)

			_, ok, err := ch.Get(worker.VideoTranscodingQueue, false)
			assert.NoError(t, err)
			assert.True(t, ok)
		})

		t.Run("should return 400 for invalid UUID", func(t *testing.T) {
			res, err := testutils.Send(r, http.MethodPost, "/medias/uuid/upload/video", nil)
			assert.NoError(t, err)

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 400, res.Result().StatusCode)
			assert.Equal(t, 400, body.Code)
			assert.EqualError(t, ErrInvalidUUID, body.Message)
		})

		t.Run("should return 404 for non-existing media", func(t *testing.T) {
			res, err := testutils.Send(r, http.MethodPost, "/medias/7b959619-7271-4fbb-a70c-b6b5b40aecaf/upload/video", nil)
			assert.NoError(t, err)

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 404, res.Result().StatusCode)
			assert.Equal(t, 404, body.Code)
			assert.Equal(t, "media could not be found", body.Message)
		})

		t.Run("should return 400 for file missing", func(t *testing.T) {})
		t.Run("should return 400 for file size too high", func(t *testing.T) {})
		t.Run("should return error for invalid media file", func(t *testing.T) {})
		t.Run("should set media status to errored on upload error", func(t *testing.T) {})
	})
}
