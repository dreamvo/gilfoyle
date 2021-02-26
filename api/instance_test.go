package api

import (
	"encoding/json"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/x/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

func TestInstance(t *testing.T) {
	s := NewServer(Options{
		Config: config.Config{
			Settings:config.SettingsConfig {
				Debug: true,
				MaxFileSize: 100,
			},
			Storage: config.StorageConfig{
				Driver: "test_fs",
			},
			Services: config.ServicesConfig{
				DB: config.DatabaseConfig{
					Dialect: "pg",
				},
			},
		},
		Logger: zap.NewExample(),
	})

	t.Run("GET /healthz", func(t *testing.T) {
		res, err := testutils.Send(s.router, http.MethodGet, "/healthz", nil)
		assert.NoError(t, err)

		var body HealthCheckResponse
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(t, 200, res.Result().StatusCode)
		assert.Equal(t, HealthCheckResponse{
			Tag:    config.Version,
			Commit: config.Commit,
			Debug: true,
			DatabaseDialect: "pg",
			MaxFileSize: 100,
			StorageDriver: "test_fs",
		}, body)
	})

	t.Run("GET /metricsz", func(t *testing.T) {
		res, err := testutils.Send(s.router, http.MethodGet, "/metricsz", nil)
		assert.NoError(t, err)

		assert.Equal(t, 200, res.Result().StatusCode)
		assert.Equal(t, "text/plain; version=0.0.4; charset=utf-8", res.Header().Get("Content-Type"))
	})

	t.Run("GET /404notfound", func(t *testing.T) {
		res, err := testutils.Send(s.router, http.MethodGet, "/404notfound", nil)
		assert.NoError(t, err)

		var body util.ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(t, 404, res.Result().StatusCode)
		assert.Equal(t, util.ErrorResponse{
			Code:    404,
			Message: "resource not found",
		}, body)
	})
}
