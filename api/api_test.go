package api

import (
	"bytes"
	"encoding/json"
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
	r = gin.Default()
	r = RegisterRoutes(r)

	t.Run("GET /health", func(t *testing.T) {
		res, err := performRequest(r, "GET", "/health", nil)
		assert.NoError(err)

		var body struct {
			Code int               `json:"code"`
			Data map[string]string `json:"data,omitempty"`
		}
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(200, body.Code)
		assert.Equal(config.Version, body.Data["tag"])
		assert.Equal(config.Commit, body.Data["commit"])
	})
}
