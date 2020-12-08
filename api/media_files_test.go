package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/media"
	"github.com/dreamvo/gilfoyle/ent/mediafile"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/storage"
	_ "github.com/mattn/go-sqlite3"
	assertTest "github.com/stretchr/testify/assert"
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
	assert := assertTest.New(t)
	r = NewServer()

	_, err := gilfoyle.NewConfig()
	if err != nil {
		panic(err)
	}

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer removeDir(gilfoyle.Config.Storage.Filesystem.DataPath)

	_, err = gilfoyle.NewStorage(storage.Filesystem)
	if err != nil {
		panic(err)
	}

	t.Run("POST /medias/:id/upload/video", func(t *testing.T) {
		t.Run("should upload file and return probe", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			m, _ := db.Client.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)

			filePath := "./__mocks__/SampleVideo_1280x720_1mb.mp4"

			file, err := os.Open(filePath)
			assert.NoError(err)
			defer file.Close()

			part1, err := writer.CreateFormFile("file", filepath.Base(filePath))
			assert.NoError(err)

			_, err = io.Copy(part1, file)
			assert.NoError(err)

			err = writer.Close()
			assert.NoError(err)

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/medias/%s/upload/video", m.ID), payload)
			assert.NoError(err)

			req.Header.Add("Content-Type", "multipart/form-data")
			req.Header.Set("Content-Type", writer.FormDataContentType())

			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			var body util.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			stat, err := os.Stat(filepath.Join("./data", m.ID.String(), "original"))
			assert.NoError(err)
			assert.Equal(int64(1055736), stat.Size())

			assert.Equal(200, res.Result().StatusCode)
			assert.Equal(map[string]interface{}{
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

			m, _ = db.Client.Media.Get(context.Background(), m.ID)

			assert.Equal(media.StatusProcessing, m.Status)

			mediaFile, _ := db.Client.MediaFile.
				Query().
				Where(mediafile.MediaTypeEQ(schema.MediaFileTypeVideo)).
				Only(context.Background())

			assert.Equal(int8(25), mediaFile.Framerate)
			assert.Equal(5.312, mediaFile.DurationSeconds)
			assert.Equal(int16(1280), mediaFile.ScaledWidth)
			assert.Equal(mediafile.EncoderPreset(schema.MediaFileEncoderPresetSource), mediaFile.EncoderPreset)
			assert.Equal(int64(1205959), mediaFile.VideoBitrate)
			assert.Equal(mediafile.MediaType(schema.MediaFileTypeVideo), mediaFile.MediaType)
		})

		t.Run("should return 400 for invalid UUID", func(t *testing.T) {
			res, err := performRequest(r, http.MethodPost, "/medias/uuid/upload/video", nil)
			assert.NoError(err)

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode)
			assert.Equal(400, body.Code)
			assert.EqualError(ErrInvalidUUID, body.Message)
		})

		t.Run("should return 404 for non-existing media", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, http.MethodPost, "/medias/7b959619-7271-4fbb-a70c-b6b5b40aecaf/upload/video", nil)
			assert.NoError(err)

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(404, res.Result().StatusCode)
			assert.Equal(404, body.Code)
			assert.Equal("media could not be found", body.Message)
		})

		t.Run("should return 400 for file missing", func(t *testing.T) {})
		t.Run("should return 400 for file size too high", func(t *testing.T) {})
		t.Run("should return error for invalid media file", func(t *testing.T) {})
		t.Run("should set media status to errored on upload error", func(t *testing.T) {})
	})
}
