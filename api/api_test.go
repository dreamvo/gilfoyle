package api

import (
	"github.com/gin-gonic/gin"
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
	r = RegisterRoutes(r, false)

	t.Run("GET /health", func(t *testing.T) {
		res, err := performRequest(r, "GET", "/health")

		assert.Equal(err, nil, "should be equal")
		assert.Equal(res.Result().StatusCode, 200, "should be equal")
	})
}
