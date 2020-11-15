package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/config"
	assertTest "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	assert := assertTest.New(t)

	t.Run("should set default values", func(t *testing.T) {
		_, err := NewConfig()
		assert.Nil(err)

		assert.Equal(&config.Config{
			Services: config.ServicesConfig{
				DB: config.DatabaseConfig{
					Dialect:  "postgres",
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "",
					Database: "gilfoyle",
				},
				Redis: config.RedisConfig{
					Host:     "localhost",
					Database: "0",
					Port:     "6379",
					Password: "",
				},
			},
			Settings: config.SettingsConfig{
				ExposeSwaggerUI: true,
				MaxFileSize:     524288000,
				Debug:           false,
			},
			Storage: config.StorageConfig{
				Class: "fs",
				Filesystem: config.FileSystemConfig{
					DataPath: "/data",
				},
				S3: config.S3Config{
					Hostname:        "",
					Port:            "",
					AccessKeyID:     "",
					SecretAccessKey: "",
					Region:          "",
					Bucket:          "",
					EnableSSL:       true,
					UsePathStyle:    false,
				},
				GCS: config.GCSConfig{
					CredentialsFile: "",
					Bucket:          ""},
				IPFS: config.IPFSConfig{
					Gateway: "gateway.ipfs.io"},
			},
		}, &Config, "should be equal")
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

		_, err := NewConfig()
		assert.Nil(err)

		assert.Equal("postgres", Config.Services.DB.Dialect)
		assert.Equal("host", Config.Services.DB.Host)
		assert.Equal("port", Config.Services.DB.Port)
		assert.Equal("user", Config.Services.DB.User)
		assert.Equal("database", Config.Services.DB.Database)
		assert.Equal("password", Config.Services.DB.Password)

		assert.Equal("redis_db", Config.Services.Redis.Database)
		assert.Equal("redis_pass", Config.Services.Redis.Password)
		assert.Equal("redis_port", Config.Services.Redis.Port)
		assert.Equal("redis_host", Config.Services.Redis.Host)

		assert.Equal("ipfs_gateway", Config.Storage.IPFS.Gateway)
	})

	t.Run("should not return error on bad file path", func(t *testing.T) {
		_, err := NewConfig("/path/to/file.wat")
		assert.Nil(err)
	})
}
