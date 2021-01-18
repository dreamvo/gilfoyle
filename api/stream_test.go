package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/x/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strings"
	"testing"
)

func TestStream(t *testing.T) {
	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() { _ = dbClient.Close() }()

	_, err := gilfoyle.NewConfig()
	if err != nil {
		t.Error(err)
	}

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer removeDir(t, gilfoyle.Config.Storage.Filesystem.DataPath)

	storageDriver, err := gilfoyle.NewStorage(storage.Filesystem)
	if err != nil {
		t.Error(err)
	}

	s := NewServer(Options{
		Logger:   zap.NewExample(),
		Storage:  storageDriver,
		Database: dbClient,
	})

	t.Run("GET /medias/{media_id}/stream/*filename", func(t *testing.T) {
		t.Run("should return the requested file", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusReady).
				Save(context.Background())
			assert.NoError(t, err)

			err = s.storage.Save(context.Background(), strings.NewReader("test"), path.Join(m.ID.String(), "low", "index.m3u8"))
			assert.NoError(t, err)

			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/low/index.m3u8", m.ID.String()), nil)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
			assert.Equal(t, "test", res.Body.String())
			assert.Equal(t, "application/x-mpegURL", res.Header().Get("Content-Type"))
			assert.Equal(t, "attachment; filename=\"index.m3u8\"", res.Header().Get("Content-Disposition"))
		})

		t.Run("should return the requested file (2)", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusReady).
				Save(context.Background())
			assert.NoError(t, err)

			err = s.storage.Save(context.Background(), strings.NewReader("test"), path.Join(m.ID.String(), "low", "000.ts"))
			assert.NoError(t, err)

			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/low/000.ts", m.ID.String()), nil)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
			assert.Equal(t, "test", res.Body.String())
			assert.Equal(t, "video/MP2T", res.Header().Get("Content-Type"))
			assert.Equal(t, "attachment; filename=\"000.ts\"", res.Header().Get("Content-Disposition"))
		})

		t.Run("should return error because of non-Ready status", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())
			assert.NoError(t, err)

			var body util.ErrorResponse
			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/low/index.m3u8", m.ID.String()), nil)
			assert.NoError(t, err)

			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, http.StatusTooEarly, res.Result().StatusCode)
			assert.Equal(t, util.ErrorResponse{
				Code:    425,
				Message: "media is not ready yet for streaming",
			}, body)
		})

		t.Run("should return error because of non-existing media", func(t *testing.T) {
			var body util.ErrorResponse
			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/low/index.m3u8", uuid.New().String()), nil)
			assert.NoError(t, err)

			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
			assert.Equal(t, util.ErrorResponse{
				Code:    404,
				Message: ErrResourceNotFound.Error(),
			}, body)
		})

		t.Run("should return error because of bad uuid", func(t *testing.T) {
			var body util.ErrorResponse
			res, err := testutils.Send(s.router, http.MethodGet, "/medias/uuid/stream/low/index.m3u8", nil)
			assert.NoError(t, err)

			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
			assert.Equal(t, util.ErrorResponse{
				Code:    400,
				Message: ErrInvalidUUID.Error(),
			}, body)
		})

		t.Run("should return non existing file", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusReady).
				Save(context.Background())
			assert.NoError(t, err)

			var body util.ErrorResponse
			res, err := testutils.Send(s.router, http.MethodGet, fmt.Sprintf("/medias/%s/stream/low/index.m3u8", m.ID.String()), nil)
			assert.NoError(t, err)

			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
			assert.Equal(t, util.ErrorResponse{
				Code:    404,
				Message: "file does not exist",
			}, body)
		})
	})
}
