package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("should set default values", func(t *testing.T) {
		cfg, err := NewConfig()
		assert.Nil(t, err)

		assert.Equal(t, &Config{
			Services: ServicesConfig{
				DB: DatabaseConfig{
					Dialect:  "postgres",
					Host:     "localhost",
					Port:     "5432",
					User:     "postgres",
					Password: "",
					Database: "gilfoyle",
				},
				RabbitMQ: RabbitMQConfig{
					Host:     "localhost",
					Username: "guest",
					Port:     5672,
					Password: "guest",
				},
			},
			Settings: SettingsConfig{
				ExposeSwaggerUI: true,
				MaxFileSize:     524288000,
				Debug:           false,
				Worker: WorkerSettings{
					Concurrency: 10,
				},
			},
			Storage: StorageConfig{
				Driver: "fs",
				Filesystem: FileSystemConfig{
					DataPath: "/data",
				},
				S3: S3Config{
					Hostname:        "",
					Port:            "",
					AccessKeyID:     "",
					SecretAccessKey: "",
					Region:          "",
					Bucket:          "",
					EnableSSL:       true,
					UsePathStyle:    false,
				},
				GCS: GCSConfig{
					CredentialsFile: "",
					Bucket:          "",
				},
			},
		}, cfg, "should be equal")
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

		cfg, err := NewConfig()
		assert.Nil(t, err)

		assert.Equal(t, "postgres", cfg.Services.DB.Dialect)
		assert.Equal(t, "host", cfg.Services.DB.Host)
		assert.Equal(t, "port", cfg.Services.DB.Port)
		assert.Equal(t, "user", cfg.Services.DB.User)
		assert.Equal(t, "database", cfg.Services.DB.Database)
		assert.Equal(t, "password", cfg.Services.DB.Password)

		assert.Equal(t, "rabbitmq_user", cfg.Services.RabbitMQ.Username)
		assert.Equal(t, "rabbitmq_pass", cfg.Services.RabbitMQ.Password)
		assert.Equal(t, 5555, cfg.Services.RabbitMQ.Port)
		assert.Equal(t, "rabbitmq_host", cfg.Services.RabbitMQ.Host)
	})

	t.Run("should not return error on bad file path", func(t *testing.T) {
		_, err := NewConfig("/path/to/file.wat")
		assert.Nil(t, err)
	})
}
