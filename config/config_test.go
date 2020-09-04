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

		assert.Equal(&Config{
			Services: servicesConfig{
				IPFS: ipfsConfig{
					Gateway: "gateway.ipfs.io",
				},
				DB: dbConfig{
					Dialect:  "postgres",
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "",
					Database: "gilfoyle",
				},
				Redis: redisConfig{
					Host:     "localhost",
					Port:     "6379",
					Database: "0",
					Password: "",
				},
			},
			Settings: settingsConfig{
				ServeDocs:   true,
				MaxFileSize: "50mb",
			},
		}, GetConfig(), "should be equal")
	})

	t.Run("should set values from env vars", func(t *testing.T) {
		defer os.Clearenv()

		_ = os.Setenv("DB_HOST", "host")
		_ = os.Setenv("DB_PORT", "port")
		_ = os.Setenv("DB_USER", "user")
		_ = os.Setenv("DB_PASSWORD", "password")
		_ = os.Setenv("DB_NAME", "database")

		_ = os.Setenv("REDIS_HOST", "redis_host")
		_ = os.Setenv("REDIS_DB", "redis_db")
		_ = os.Setenv("REDIS_PORT", "redis_port")
		_ = os.Setenv("REDIS_PASSWORD", "redis_pass")

		_ = os.Setenv("IPFS_GATEWAY", "ipfs_gateway")

		err := New()
		assert.Nil(err)

		assert.Equal("postgres", GetConfig().Services.DB.Dialect)
		assert.Equal("host", GetConfig().Services.DB.Host)
		assert.Equal("port", GetConfig().Services.DB.Port)
		assert.Equal("user", GetConfig().Services.DB.User)
		assert.Equal("database", GetConfig().Services.DB.Database)
		assert.Equal("password", GetConfig().Services.DB.Password)

		assert.Equal("redis_db", GetConfig().Services.Redis.Database)
		assert.Equal("redis_pass", GetConfig().Services.Redis.Password)
		assert.Equal("redis_port", GetConfig().Services.Redis.Port)
		assert.Equal("redis_host", GetConfig().Services.Redis.Host)

		assert.Equal("ipfs_gateway", GetConfig().Services.IPFS.Gateway)
	})

	t.Run("should not return error on bad file path", func(t *testing.T) {
		err := New("/path/to/file.wat")
		assert.Nil(err)
	})
}
