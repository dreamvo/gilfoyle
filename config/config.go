//go:generate sh -c "go run ../cmd/main.go config | tee ../.support/config/defaults.yml"
package config

import (
	"github.com/jinzhu/configor"
	"time"
)

const (
	OriginalPolicyDelete = "Delete"
	OriginalPolicyRetain = "Retain"
)

// OriginalPolicy defines how to manage original files after processing succeeded
type OriginalPolicy string

// StorageDriver is a kind of storage backend
type StorageDriver string

// Config defines the application's settings
type Config struct {
	Services ServicesConfig `yaml:"services" json:"services"`
	Settings SettingsConfig `yaml:"settings" json:"settings"`
	Storage  StorageConfig  `yaml:"storage" json:"storage"`
}

type ServicesConfig struct {
	DB       DatabaseConfig `yaml:"db" json:"db"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq" json:"rabbitmq"`
}

type SettingsConfig struct {
	ExposeSwaggerUI bool            `yaml:"expose_swagger_ui" json:"expose_swagger_ui" default:"true"`
	MaxFileSize     int64           `yaml:"max_file_size" json:"max_file_size" default:"524288000"`
	Debug           bool            `yaml:"debug" json:"debug" default:"false" env:"APP_DEBUG"`
	Worker          WorkerSettings  `yaml:"worker" json:"worker"`
	Encoding        EncoderSettings `yaml:"encoding" json:"encoding"`
}

type StorageConfig struct {
	Driver     string           `yaml:"driver" json:"driver" default:"fs" env:"STORAGE_DRIVER"`
	Filesystem FileSystemConfig `yaml:"fs" json:"fs"`
	S3         S3Config         `yaml:"s3" json:"s3"`
	GCS        GCSConfig        `yaml:"gcs" json:"gcs"`
}

type FileSystemConfig struct {
	DataPath string `yaml:"data_path" json:"data_path" default:"/data" env:"FS_DATA_PATH"`
}

type S3Config struct {
	Hostname        string `yaml:"hostname" json:"hostname" default:"" env:"S3_HOSTNAME"`
	Port            string `yaml:"port" json:"port" default:"" env:"S3_PORT"`
	AccessKeyID     string `yaml:"access_key_id" json:"access_key_id" env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `yaml:"secret_access_key" json:"secret_access_key" env:"S3_SECRET_ACCESS_KEY"`
	Region          string `yaml:"region" json:"region" env:"S3_REGION"`
	Bucket          string `yaml:"bucket" json:"bucket" env:"S3_BUCKET"`
	EnableSSL       bool   `yaml:"enable_ssl" json:"enable_ssl" default:"true" env:"S3_ENABLE_SSL"`
	UsePathStyle    bool   `yaml:"use_path_style" json:"use_path_style" default:"false" env:"S3_ENABLE_PATH_STYLE"`
}

type GCSConfig struct {
	CredentialsFile string `yaml:"credentials_file" json:"credentials_file" default:"" env:"GCS_CREDENTIALS_FILE"`
	Bucket          string `yaml:"bucket" json:"bucket" default:"" env:"GCS_BUCKET"`
}

type DatabaseConfig struct {
	Dialect  string `yaml:"-" json:"-" default:"postgres"`
	Host     string `yaml:"host" json:"host" default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" json:"port" default:"5432" env:"DB_PORT"`
	User     string `yaml:"user" json:"user" default:"postgres" env:"DB_USER"`
	Password string `yaml:"password" json:"password" default:"" env:"DB_PASSWORD"`
	Database string `yaml:"db_name" json:"db_name" default:"gilfoyle" env:"DB_NAME"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host" json:"host" default:"localhost" env:"RABBITMQ_HOST"`
	Port     int    `yaml:"port" json:"port" default:"5672" env:"RABBITMQ_PORT"`
	Username string `yaml:"username" json:"username" default:"guest" env:"RABBITMQ_USER"`
	Password string `yaml:"password" json:"password" default:"guest" env:"RABBITMQ_PASSWORD"`
}

type WorkerSettings struct {
	Concurrency uint `yaml:"concurrency" json:"concurrency" default:"10" env:"WORKER_CONCURRENCY"`
}

type EncoderSettings struct {
	OriginalPolicy string      `yaml:"original_policy" json:"original_policy"`
	Renditions     []Rendition `yaml:"renditions" json:"renditions"`
}

type Rendition struct {
	Name         string `yaml:"name" json:"name" default:"default"`
	Width        int    `yaml:"width" json:"width" default:"842"`
	Height       int    `yaml:"height" json:"height" default:"480"`
	VideoBitrate int    `yaml:"video_bitrate" json:"video_bitrate" default:"1400000"`
	AudioBitrate int    `yaml:"audio_bitrate" json:"audio_bitrate" default:"128000"`
	Framerate    int    `yaml:"framerate" json:"framerate" default:"0"`
	VideoCodec   string `yaml:"video_codec" json:"video_codec" default:"h264"`
	AudioCodec   string `yaml:"audio_codec" json:"audio_codec" default:"aac"`
}

// NewConfig creates a new config object
// and load values from environment variables or config file.
// File paths can be both relative and absolute.
func NewConfig(files ...string) (*Config, error) {
	var config Config
	err := configor.New(&configor.Config{
		AutoReload:           false,
		AutoReloadInterval:   30 * time.Second,
		Debug:                false,
		Silent:               false,
		Verbose:              false,
		ErrorOnUnmatchedKeys: true,
	}).Load(&config, files...)
	return &config, err
}
