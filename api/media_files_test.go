package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	assertTest "github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMediaFiles(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.New()
	r = RegisterRoutes(r)

	t.Run("POST /medias/{id}/upload", func(t *testing.T) {
		t.Run("(WIP) should return 200", func(t *testing.T) {
			_, err := performRequest(r, http.MethodPost, "/medias/uuid/upload", nil)
			assert.NoError(err, "should be equal")

			//assert.Equal(200, res.Result().StatusCode, "should be equal")
		})
	})
}
