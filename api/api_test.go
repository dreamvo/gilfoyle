package api

import (
	"bytes"
	"encoding/json"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/gin-gonic/gin"
	assertTest "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var r *gin.Engine

func performRequest(r http.Handler, method, path string, body interface{}) (*httptest.ResponseRecorder, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, path, bytes.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}

func TestApi(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.New()
	r = RegisterMiddlewares(r)
	r = RegisterRoutes(r)

	t.Run("GET /healthz", func(t *testing.T) {
		res, err := performRequest(r, "GET", "/healthz", nil)
		assert.NoError(err)

		var body HealthCheckResponse
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(200, res.Result().StatusCode)
		assert.Equal(config.Version, body.Tag)
		assert.Equal(config.Commit, body.Commit)
	})

	t.Run("GET /404notfound", func(t *testing.T) {
		res, err := performRequest(r, "GET", "/404notfound", nil)
		assert.NoError(err)

		var body util.ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(404, res.Result().StatusCode)
		assert.Equal(404, body.Code)
		assert.Equal("resource not found", body.Message)
	})
}
