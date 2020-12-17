package api

import (
	"encoding/json"
	"github.com/dreamvo/gilfoyle/api/util"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/x/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var r *gin.Engine

func TestApi(t *testing.T) {
	s := NewServer(Options{})
	r = s.router

	t.Run("GET /healthz", func(t *testing.T) {
		res, err := testutils.Send(r, http.MethodGet, "/healthz", nil)
		assert.NoError(t, err)

		var body HealthCheckResponse
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(t, 200, res.Result().StatusCode)
		assert.Equal(t, config.Version, body.Tag)
		assert.Equal(t, config.Commit, body.Commit)
	})

	t.Run("GET /404notfound", func(t *testing.T) {
		res, err := testutils.Send(r, http.MethodGet, "/404notfound", nil)
		assert.NoError(t, err)

		var body util.ErrorResponse
		_ = json.NewDecoder(res.Body).Decode(&body)

		assert.Equal(t, 404, res.Result().StatusCode)
		assert.Equal(t, 404, body.Code)
		assert.Equal(t, "resource not found", body.Message)
	})
}
