package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/x/testutils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestMedias(t *testing.T) {
	t.Run("GET /medias", func(t *testing.T) {
		t.Run("should return empty array", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			res, err := testutils.Send(s.router, http.MethodGet, "/medias", nil)
			assert.NoError(t, err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, res.Result().StatusCode, 200, "should be equal")
			assert.Equal(t, 200, body.Code)
			assert.Equal(t, []ent.Media{}, body.Data)
		})

		t.Run("should return latest medias", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			for i := 0; i < 5; i++ {
				_, _ = dbClient.Media.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.MediaStatusAwaitingUpload).
					Save(context.Background())
			}

			res, err := testutils.Send(s.router, http.MethodGet, "/medias", nil)
			assert.NoError(t, err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, res.Result().StatusCode, 200, "should be equal")
			assert.Equal(t, 200, body.Code)
			assert.Equal(t, 5, len(body.Data))
		})

		t.Run("should limit results to 2", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			for i := 0; i < 3; i++ {
				_, _ = dbClient.Media.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.MediaStatusAwaitingUpload).
					Save(context.Background())
			}

			res, err := testutils.Send(s.router, http.MethodGet, "/medias?limit=2", nil)
			assert.NoError(t, err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, res.Result().StatusCode, 200, "should be equal")
			assert.Equal(t, 200, body.Code)
			assert.Equal(t, 2, len(body.Data))
		})

		t.Run("should return results with offset 1", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			v, _ := dbClient.Media.
				Create().
				SetTitle("video1").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())

			_, _ = dbClient.Media.
				Create().
				SetTitle("video2").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())

			res, err := testutils.Send(s.router, http.MethodGet, "/medias?offset=1", nil)
			assert.NoError(t, err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, res.Result().StatusCode, 200, "should be equal")
			assert.Equal(t, 200, body.Code)
			assert.Equal(t, 1, len(body.Data))
			assert.Equal(t, v.ID.String(), body.Data[0].ID.String())
		})
	})

	t.Run("GET /medias/:id", func(t *testing.T) {
		t.Run("should return error for invalid UUID", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			res, err := testutils.Send(s.router, http.MethodGet, "/medias/uuid", nil)
			assert.NoError(t, err, "should be equal")

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 400, res.Result().StatusCode, "should be equal")
			assert.Equal(t, 400, body.Code)
			assert.Equal(t, "invalid UUID provided", body.Message)
		})

		t.Run("should return media", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			v, _ := dbClient.Media.
				Create().
				SetTitle("no u").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())

			res, err := testutils.Send(s.router, http.MethodGet, "/medias/"+v.ID.String(), nil)
			assert.NoError(t, err, "should be equal")

			var body struct {
				Code int       `json:"code"`
				Data ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 200, res.Result().StatusCode, "should be equal")
			assert.Equal(t, 200, body.Code)
			assert.Equal(t, v.Title, body.Data.Title)
		})
	})

	t.Run("DELETE /medias/:id", func(t *testing.T) {
		t.Run("should delete newly created media", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			v, _ := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())

			res, err := testutils.Send(s.router, http.MethodDelete, "/medias/"+v.ID.String(), nil)
			assert.NoError(t, err, "should be equal")

			assert.Equal(t, res.Result().StatusCode, 200, "should be equal")

			res, err = testutils.Send(s.router, http.MethodDelete, "/medias/"+v.ID.String(), nil)
			assert.NoError(t, err, "should be equal")

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 404, res.Code)
			assert.Equal(t, "resource not found", body.Message)
		})

		t.Run("should return error on invalid uid", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			res, err := testutils.Send(s.router, http.MethodDelete, "/medias/uuid", nil)
			assert.NoError(t, err, "should be equal")

			var body util.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, res.Result().StatusCode, 400, "should be equal")
			assert.Equal(t, 400, body.Code)
			assert.Equal(t, "invalid UUID provided", body.Message)
		})
	})

	t.Run("POST /medias", func(t *testing.T) {
		t.Run("should create a new media", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			res, err := testutils.Send(s.router, http.MethodPost, "/medias", CreateMedia{
				Title: "test",
			})
			assert.NoError(t, err)

			var body util.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 200, res.Result().StatusCode)
			assert.Equal(t, "test", body.Data.(map[string]interface{})["title"])
			assert.Equal(t, "AwaitingUpload", body.Data.(map[string]interface{})["status"])
		})

		t.Run("should return validation error (1)", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			res, err := testutils.Send(s.router, http.MethodPost, "/medias", CreateMedia{
				Title: "Vitae sunt aspernatur quia sunt blanditiis at et excepturi. Doloribus non ut minus saepe. Quas enim minus modi possimus. Blanditiis eius in ipsam incidunt rem et. Rerum blanditiis consequatur facilis eos quia. Sed autem inventore iure ducimus voluptas voluptas.",
			})
			assert.NoError(t, err, "should be equal")

			var body util.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 400, res.Result().StatusCode, "should be equal")
			assert.Equal(t, "Some parameters are missing or invalid", body.Message)
			assert.Equal(t, map[string]string{
				"title": "Title must be at maximum 255 characters in length",
			}, body.Fields)
		})

		t.Run("should return validation error (2)", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			res, err := testutils.Send(s.router, http.MethodPost, "/medias", nil)
			assert.NoError(t, err, "should be equal")

			var body util.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 400, res.Result().StatusCode, "should be equal")
			assert.Equal(t, "Bad request", body.Message)
			assert.Equal(t, map[string]string(nil), body.Fields)
		})
	})

	t.Run("PATCH /medias/:id", func(t *testing.T) {
		t.Run("should update a media", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())
			assert.NoError(t, err)

			res, err := testutils.Send(s.router, http.MethodPatch, "/medias/"+m.ID.String(), CreateMedia{
				Title: "test2",
			})
			assert.NoError(t, err)

			var body util.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 200, res.Result().StatusCode)
			assert.Equal(t, "test2", body.Data.(map[string]interface{})["title"])
			assert.Equal(t, "AwaitingUpload", body.Data.(map[string]interface{})["status"])
		})

		t.Run("should return validation error", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			s := NewServer(Options{
				Database: dbClient,
				Logger:   zap.NewExample(),
			})

			m, err := dbClient.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusAwaitingUpload).
				Save(context.Background())
			assert.NoError(t, err)

			res, err := testutils.Send(s.router, http.MethodPatch, "/medias/"+m.ID.String(), CreateMedia{
				Title: "Vitae sunt aspernatur quia sunt blanditiis at et excepturi. Doloribus non ut minus saepe. Quas enim minus modi possimus. Blanditiis eius in ipsam incidunt rem et. Rerum blanditiis consequatur facilis eos quia. Sed autem inventore iure ducimus voluptas voluptas.",
			})
			assert.NoError(t, err, "should be equal")

			var body util.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(t, 400, res.Result().StatusCode, "should be equal")
			assert.Equal(t, "Some parameters are missing or invalid", body.Message)
			assert.Equal(t, map[string]string{
				"title": "Title must be at maximum 255 characters in length",
			}, body.Fields)
		})

		t.Run("should return validation error because of bad UUID", func(t *testing.T) {})
		t.Run("should return resource not found", func(t *testing.T) {})
	})
}
