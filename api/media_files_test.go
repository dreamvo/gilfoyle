package api

import (
	"github.com/dreamvo/gilfoyle/api/db"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	assertTest "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestMediaFiles(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.Default()
	r = RegisterRoutes(r, RouterOptions{
		ExposeSwaggerUI: false,
	})

	t.Run("POST /medias/{id}/upload", func(t *testing.T) {
		t.Run("(WIP) should return 200", func(t *testing.T) {
			res, err := performRequest(r, "POST", "/medias/uuid/upload", nil)
			assert.NoError(err, "should be equal")

			assert.Equal(200, res.Result().StatusCode, "should be equal")
		})

		t.Run("should return error on invalid uuid", func(t *testing.T) {
			db.Client = enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer db.Client.Close()

			res, err := performRequest(r, "DELETE", "/medias/uuid", nil)
			assert.Equal(nil, err, "should be equal")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(400, res.Result().StatusCode, "should be equal")
			assert.JSONEq("{\"code\": 400, \"message\":\"invalid UUID provided\"}", string(body))
		})
	})
}
