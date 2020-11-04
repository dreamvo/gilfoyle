package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/httputils"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	assertTest "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMedias(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.Default()
	r = RegisterRoutes(r)

	t.Run("GET /medias", func(t *testing.T) {
		t.Run("should return empty array", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, http.MethodGet, "/medias", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal([]ent.Media{}, body.Data)
		})

		t.Run("should return latest medias", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			for i := 0; i < 5; i++ {
				_, _ = db.Client.Media.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.MediaStatusProcessing).
					Save(context.Background())
			}

			res, err := performRequest(r, http.MethodGet, "/medias", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(5, len(body.Data))
		})

		t.Run("should limit results to 2", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			for i := 0; i < 3; i++ {
				_, _ = db.Client.Media.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.MediaStatusProcessing).
					Save(context.Background())
			}

			res, err := performRequest(r, http.MethodGet, "/medias?limit=2", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(2, len(body.Data))
		})

		t.Run("should return results with offset 1", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Media.
				Create().
				SetTitle("video1").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())

			_, _ = db.Client.Media.
				Create().
				SetTitle("video2").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, http.MethodGet, "/medias?offset=1", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(1, len(body.Data))
			assert.Equal(v.ID.String(), body.Data[0].ID.String())
		})
	})

	t.Run("GET /medias/{id}", func(t *testing.T) {
		t.Run("should return error for invalid UUID", func(t *testing.T) {
			res, err := performRequest(r, http.MethodGet, "/medias/uuid", nil)
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal(400, body.Code)
			assert.Equal("invalid UUID provided", body.Message)
		})

		t.Run("should return media", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Media.
				Create().
				SetTitle("no u").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, http.MethodGet, "/medias/"+v.ID.String(), nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int       `json:"code"`
				Data ent.Media `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(200, res.Result().StatusCode, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(v.Title, body.Data.Title)
		})
	})

	t.Run("DELETE /medias/{id}", func(t *testing.T) {
		t.Run("should delete newly created media", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, http.MethodDelete, "/medias/"+v.ID.String(), nil)
			assert.NoError(err, "should be equal")

			assert.Equal(res.Result().StatusCode, 200, "should be equal")

			res, err = performRequest(r, http.MethodDelete, "/medias/"+v.ID.String(), nil)
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(404, res.Code)
			assert.Equal("resource not found", body.Message)
		})

		t.Run("should return error on invalid uid", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, http.MethodDelete, "/medias/uuid", nil)
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 400, "should be equal")
			assert.Equal(400, body.Code)
			assert.Equal("invalid UUID provided", body.Message)
		})
	})

	t.Run("POST /medias", func(t *testing.T) {
		t.Run("should create a new media", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, http.MethodPost, "/medias", CreateMedia{
				Title: "test",
			})
			assert.NoError(err)

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(200, res.Result().StatusCode)
			assert.Equal("test", body.Data.(map[string]interface{})["title"])
			assert.Equal("processing", body.Data.(map[string]interface{})["status"])
		})

		t.Run("should return validation error (1)", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, http.MethodPost, "/medias", CreateMedia{
				Title: "Vitae sunt aspernatur quia sunt blanditiis at et excepturi. Doloribus non ut minus saepe. Quas enim minus modi possimus. Blanditiis eius in ipsam incidunt rem et. Rerum blanditiis consequatur facilis eos quia. Sed autem inventore iure ducimus voluptas voluptas.",
			})
			assert.NoError(err, "should be equal")

			var body httputils.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal("Some parameters are missing or invalid", body.Message)
			assert.Equal(map[string]string{
				"title": "Title must be at maximum 255 characters in length",
			}, body.Fields)
		})

		t.Run("should return validation error (2)", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, http.MethodPost, "/medias", nil)
			assert.NoError(err, "should be equal")

			var body httputils.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal("Bad request", body.Message)
			assert.Equal(map[string]string(nil), body.Fields)
		})
	})

	t.Run("PATCH /medias/{id}", func(t *testing.T) {
		t.Run("should update a media", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			m, err := db.Client.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())
			assert.NoError(err)

			res, err := performRequest(r, http.MethodPatch, "/medias/"+m.ID.String(), CreateMedia{
				Title: "test2",
			})
			assert.NoError(err)

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(200, res.Result().StatusCode)
			assert.Equal("test2", body.Data.(map[string]interface{})["title"])
			assert.Equal("processing", body.Data.(map[string]interface{})["status"])
		})

		t.Run("should return validation error", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			m, err := db.Client.Media.
				Create().
				SetTitle("test").
				SetStatus(schema.MediaStatusProcessing).
				Save(context.Background())
			assert.NoError(err)

			res, err := performRequest(r, http.MethodPatch, "/medias/"+m.ID.String(), CreateMedia{
				Title: "Vitae sunt aspernatur quia sunt blanditiis at et excepturi. Doloribus non ut minus saepe. Quas enim minus modi possimus. Blanditiis eius in ipsam incidunt rem et. Rerum blanditiis consequatur facilis eos quia. Sed autem inventore iure ducimus voluptas voluptas.",
			})
			assert.NoError(err, "should be equal")

			var body httputils.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal("Some parameters are missing or invalid", body.Message)
			assert.Equal(map[string]string{
				"title": "Title must be at maximum 255 characters in length",
			}, body.Fields)
		})

		t.Run("should return validation error because of bad UUID", func(t *testing.T) {})
		t.Run("should return resource not found", func(t *testing.T) {})
	})
}
