package dashboard

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

var server *Server

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

func TestServer(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err)
	}
	server, err = NewServer(logger, "http://localhost:1234")
	if err != nil {
		t.Error(err)
	}
	defer func() { _ = logger.Sync() }()

	t.Run("should redirect / to /app", func(t *testing.T) {
		res, err := performRequest(server.router, http.MethodGet, "/", nil)
		assert.NoError(t, err)

		assert.Equal(t, 307, res.Result().StatusCode)
		assert.Equal(t, "/app", res.Result().Header.Get("Location"))
	})
}
