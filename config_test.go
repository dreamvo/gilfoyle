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
				RabbitMQ: config.RabbitMQConfig{
					Host:     "localhost",
					Username: "guest",
					Port:     5672,
					Password: "guest",
				},
			},
			Settings: config.SettingsConfig{
				ExposeSwaggerUI: true,
				MaxFileSize:     524288000,
				Debug:           false,
				Worker: config.WorkerSettings{
					Concurrency: 3,
				},
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

		_ = os.Setenv("RABBITMQ_HOST", "rabbitmq_host")
		_ = os.Setenv("RABBITMQ_USER", "rabbitmq_user")
		_ = os.Setenv("RABBITMQ_PORT", "5555")
		_ = os.Setenv("RABBITMQ_PASSWORD", "rabbitmq_pass")

		_ = os.Setenv("IPFS_GATEWAY", "ipfs_gateway")

		_, err := NewConfig()
		assert.Nil(err)

		assert.Equal("postgres", Config.Services.DB.Dialect)
		assert.Equal("host", Config.Services.DB.Host)
		assert.Equal("port", Config.Services.DB.Port)
		assert.Equal("user", Config.Services.DB.User)
		assert.Equal("database", Config.Services.DB.Database)
		assert.Equal("password", Config.Services.DB.Password)

		assert.Equal("rabbitmq_user", Config.Services.RabbitMQ.Username)
		assert.Equal("rabbitmq_pass", Config.Services.RabbitMQ.Password)
		assert.Equal(5555, Config.Services.RabbitMQ.Port)
		assert.Equal("rabbitmq_host", Config.Services.RabbitMQ.Host)

		assert.Equal("ipfs_gateway", Config.Storage.IPFS.Gateway)
	})

	t.Run("should not return error on bad file path", func(t *testing.T) {
		_, err := NewConfig("/path/to/file.wat")
		assert.Nil(err)
	})
}
