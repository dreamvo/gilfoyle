package dashboard

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var server *Server

func performRequest(r http.Handler, method, path string, body interface{}, headers map[string]string) (*httptest.ResponseRecorder, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}

func TestServer(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err)
	}

	server, err = NewServer(logger, "http://gilfoyle.localhost")
	if err != nil {
		t.Error(err)
	}
	defer func() { _ = logger.Sync() }()

	t.Run("should redirect / to /app", func(t *testing.T) {
		res, err := performRequest(server.router, http.MethodGet, "/", nil, map[string]string{})
		assert.NoError(t, err)

		assert.Equal(t, 307, res.Result().StatusCode)
		assert.Equal(t, "/app", res.Result().Header.Get("Location"))
	})

	t.Run("should forward request to /healthz", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("http://gilfoyle.localhost").
			Post("/healthz").
			MatchParam("test", "1").
			MatchHeader("Content-Type", "text/html").
			MatchHeader("X-Test", "test").
			Reply(201).
			SetHeader("X-Gilfoyle-Header", "test").
			JSON(map[string]string{"foo": "bar"})

		headers := map[string]string{
			"Content-Type": "text/html",
			"X-Test":       "test",
		}

		res, err := performRequest(server.router, http.MethodPost, "/api/proxy/healthz?test=1", nil, headers)
		assert.NoError(t, err)

		var body map[string]string
		err = json.NewDecoder(res.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, 201, res.Result().StatusCode)
		assert.Equal(t, map[string]string{"foo": "bar"}, body)
		assert.Equal(t, "test", res.Header().Get("X-Gilfoyle-Header"))
	})

	t.Run("should fail to forward request to /healthz", func(t *testing.T) {
		res, err := performRequest(server.router, http.MethodPost, "/api/proxy/healthz", nil, map[string]string{})
		assert.NoError(t, err)

		assert.Equal(t, 500, res.Result().StatusCode)
	})

	t.Run("should fail to parse endpoint URL", func(t *testing.T) {
		logger, err := zap.NewProduction()
		assert.NoError(t, err)

		server, err := NewServer(logger, "\\wrongendpoint:é")
		assert.Nil(t, server)
		assert.EqualError(t, err, "parse \"\\\\wrongendpoint:é\": first path segment in URL cannot contain colon")
	})
}
