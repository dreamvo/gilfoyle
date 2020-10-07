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
	"io/ioutil"
	"testing"
)

func TestVideo(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.Default()
	r = RegisterRoutes(r, RouterOptions{
		ExposeSwaggerUI: false,
	})

	t.Run("GET /videos", func(t *testing.T) {
		t.Run("should return empty array", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "GET", "/videos", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Video `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal([]ent.Video{}, body.Data)
		})

		t.Run("should return latest videos", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			for i := 0; i < 5; i++ {
				_, _ = db.Client.Video.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.VideoStatusProcessing).
					Save(context.Background())
			}

			res, err := performRequest(r, "GET", "/videos", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Video `json:"data,omitempty"`
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
				_, _ = db.Client.Video.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.VideoStatusProcessing).
					Save(context.Background())
			}

			res, err := performRequest(r, "GET", "/videos?limit=2", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Video `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(2, len(body.Data))
		})

		t.Run("should return results with offset 1", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Video.
				Create().
				SetTitle("video1").
				SetStatus(schema.VideoStatusProcessing).
				Save(context.Background())

			_, _ = db.Client.Video.
				Create().
				SetTitle("video2").
				SetStatus(schema.VideoStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, "GET", "/videos?offset=1", nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int         `json:"code"`
				Data []ent.Video `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(1, len(body.Data))
			assert.Equal(v.ID.String(), body.Data[0].ID.String())
		})
	})

	t.Run("GET /videos/{id}", func(t *testing.T) {
		t.Run("should return error for invalid UUID", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/videos/uuid", nil)
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal(400, body.Code)
			assert.Equal("invalid UUID provided", body.Message)
		})

		t.Run("should return video", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Video.
				Create().
				SetTitle("no u").
				SetStatus(schema.VideoStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, "GET", "/videos/"+v.ID.String(), nil)
			assert.NoError(err, "should be equal")

			var body struct {
				Code int       `json:"code"`
				Data ent.Video `json:"data,omitempty"`
			}
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(200, res.Result().StatusCode, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(v.Title, body.Data.Title)
		})
	})

	t.Run("DELETE /videos/{id}", func(t *testing.T) {
		t.Run("should delete newly created video", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Video.
				Create().
				SetTitle("test").
				SetStatus(schema.VideoStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, "DELETE", "/videos/"+v.ID.String(), nil)
			assert.NoError(err, "should be equal")

			assert.Equal(res.Result().StatusCode, 200, "should be equal")

			res, err = performRequest(r, "DELETE", "/videos/"+v.ID.String(), nil)
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(404, res.Code)
			assert.Equal("resource not found", body.Message)
		})

		t.Run("should return error on invalid uid", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "DELETE", "/videos/uuid", nil)
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 400, "should be equal")
			assert.Equal(400, body.Code)
			assert.Equal("invalid UUID provided", body.Message)
		})
	})

	t.Run("POST /videos", func(t *testing.T) {
		t.Run("should create a new video", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "POST", "/videos", CreateVideo{
				Title: "test",
			})
			assert.NoError(err, "should be equal")

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(200, res.Result().StatusCode, "should be equal")
			assert.Equal("test", body.Data.(map[string]interface{})["title"])
			assert.Equal("processing", body.Data.(map[string]interface{})["status"])
		})

		t.Run("should return validation error", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "POST", "/videos", CreateVideo{
				Title: "",
			})
			assert.NoError(err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal("ent: validator failed for field \"title\": value is less than the required length", body.Message)
		})

		t.Run("should return validation error", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "POST", "/videos", nil)
			assert.NoError(err, "should be equal")

			var body httputils.ValidationErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal("Some parameters are missing or invalid", body.Message)
			assert.Equal(map[string]httputils.ValidationField{
				"Title": {
					Message: "field is invalid",
					Type:   "string",
				},
			}, body.Fields)
		})
	})

	t.Run("PATCH /videos/{id}", func(t *testing.T) {})

	t.Run("POST /videos/{id}/upload", func(t *testing.T) {
		t.Run("(WIP) should return 200", func(t *testing.T) {
			res, err := performRequest(r, "POST", "/videos/uuid/upload", nil)
			assert.NoError(err, "should be equal")

			assert.Equal(200, res.Result().StatusCode, "should be equal")
		})

		t.Run("should return error on invalid uuid", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "DELETE", "/videos/uuid", nil)
			assert.Equal(nil, err, "should be equal")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.JSONEq("{\"code\": 400, \"message\":\"invalid UUID provided\"}", string(body))
		})
	})
}
