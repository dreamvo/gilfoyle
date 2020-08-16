package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/httputils"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	assertTest "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r *gin.Engine

func performRequest(r http.Handler, method, path string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}

func TestApi(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.Default()
	RegisterRoutes(r)

	t.Run("GET /v1/videos", func(t *testing.T) {
		t.Run("should return empty array", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "GET", "/v1/videos")
			assert.Equal(nil, err, "should be equal")

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal([]interface{}{}, body.Data)
		})

		t.Run("should return latest videos", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			for i := 0; i < 10; i++ {
				_, _ = db.Client.Video.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.VideoStatusProcessing).
					Save(context.Background())
			}

			res, err := performRequest(r, "GET", "/v1/videos")
			assert.Equal(nil, err, "should be equal")

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(10, len(body.Data.([]interface{})))
		})

		t.Run("should limit results to 5", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			for i := 0; i < 10; i++ {
				_, _ = db.Client.Video.
					Create().
					SetTitle(fmt.Sprintf("%d", i)).
					SetStatus(schema.VideoStatusProcessing).
					Save(context.Background())
			}

			res, err := performRequest(r, "GET", "/v1/videos?limit=5")
			assert.Equal(nil, err, "should be equal")

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
			assert.Equal(5, len(body.Data.([]interface{})))
		})
	})

	t.Run("GET /v1/videos/{id}", func(t *testing.T) {
		t.Run("should return error for invalid UUID", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/v1/videos/uuid")
			assert.Equal(nil, err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.Equal(400, body.Code)
			assert.Equal("invalid UUID provided", body.Message)
		})
	})

	t.Run("DELETE /v1/videos/{id}", func(t *testing.T) {
		t.Run("should delete newly created video", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			v, _ := db.Client.Video.
				Create().
				SetTitle("test").
				SetStatus(schema.VideoStatusProcessing).
				Save(context.Background())

			res, err := performRequest(r, "DELETE", "/v1/videos/"+v.ID.String())
			assert.Equal(nil, err, "should be equal")

			var body httputils.DataResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(200, body.Code)
		})

		t.Run("should return error on invalid uid", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "DELETE", "/v1/videos/uuid")
			assert.Equal(nil, err, "should be equal")

			var body httputils.ErrorResponse
			_ = json.NewDecoder(res.Body).Decode(&body)

			assert.Equal(res.Result().StatusCode, 400, "should be equal")
			assert.Equal(400, body.Code)
			assert.Equal("invalid UUID provided", body.Message)
		})
	})

	t.Run("POST /v1/videos", func(t *testing.T) {})

	t.Run("PATCH /v1/videos/{id}", func(t *testing.T) {})
}
