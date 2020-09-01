package config

import (
	assertTest "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	assert := assertTest.New(t)

	t.Run("should set default values", func(t *testing.T) {
		err := New()
		assert.Nil(err)

		assert.Equal(GetConfig(), &Config{
			Services: servicesConfig{
				IPFS: ipfsConfig{
					Gateway: "gateway.ipfs.io",
				},
				DB: dbConfig{
					Dialect:  "postgres",
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "secret",
					Database: "gilfoyle",
				},
				Redis: redisConfig{
					Host:     "localhost",
					Port:     "6379",
					Password: "",
				},
			},
			Settings: settingsConfig{
				ServeDocs:   true,
				MaxFileSize: "50mb",
			},
		}, "should be equal")
	})

	t.Run("should set values from env vars", func(t *testing.T) {
		_ = os.Setenv("DB_DIALECT", "dialect")
		_ = os.Setenv("DB_HOST", "host")
		_ = os.Setenv("DB_PORT", "port")
		_ = os.Setenv("DB_USER", "user")
		_ = os.Setenv("DB_PASSWORD", "password")
		_ = os.Setenv("DB_NAME", "database")

		err := New()
		assert.Nil(err)

		assert.Equal(GetConfig().Services.DB.Dialect, "dialect")
		assert.Equal(GetConfig().Services.DB.Host, "host")
		assert.Equal(GetConfig().Services.DB.Port, "port")
		assert.Equal(GetConfig().Services.DB.User, "user")
		assert.Equal(GetConfig().Services.DB.Database, "database")
		assert.Equal(GetConfig().Services.DB.Password, "password")
	})

	t.Run("should not return error on bad file path", func(t *testing.T) {
		err := New("/path/to/file.wat")
		assert.Nil(err)
	})
}
