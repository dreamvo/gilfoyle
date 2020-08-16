package v1

import (
	"context"
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	assertTest "github.com/stretchr/testify/assert"
	"io/ioutil"
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

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.JSONEq("{\"code\": 200, \"data\": []}", string(body))
		})
	})

	t.Run("GET /v1/videos/{id}", func(t *testing.T) {
		t.Run("should return error for invalid UUID", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/v1/videos/uuid")
			assert.Equal(nil, err, "should be equal")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.JSONEq("{\"code\":400,\"message\":\"invalid UUID provided\"}", string(body))
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

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.JSONEq("{\"code\": 200}", string(body))
		})

		t.Run("should return error on invalid uid", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "DELETE", "/v1/videos/uuid")
			assert.Equal(nil, err, "should be equal")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(res.Result().StatusCode, 400, "should be equal")
			assert.JSONEq("{\"code\": 400, \"message\":\"invalid UUID provided\"}", string(body))
		})
	})
}
