package config

import (
	"github.com/jinzhu/configor"
)

// StorageClass is a choice of storage backend
type StorageClass string

// Config defines the application's settings
type Config struct {
	Services servicesConfig `yaml:"services" json:"services"`
	Settings settingsConfig `yaml:"settings" json:"settings"`
	Storage  storageConfig  `yaml:"storage" json:"storage"`
}

type servicesConfig struct {
	DB    dbConfig    `yaml:"db" json:"db"`
	Redis redisConfig `yaml:"redis" json:"redis"`
}

type settingsConfig struct {
	ExposeSwaggerUI bool   `yaml:"expose_swagger_ui" json:"expose_swagger_ui" default:"true"`
	MaxFileSize     string `yaml:"max_file_size" json:"max_file_size" default:"50Mi"`
}

type storageConfig struct {
	Class      StorageClass     `yaml:"class" json:"class" default:"fs" env:"STORAGE_CLASS"`
	CachePath  string           `yaml:"cache_path" json:"cache_path" default:"/tmp" env:"CACHE_PATH"`
	Filesystem fileSystemConfig `yaml:"fs" json:"fs"`
	S3         s3Config         `yaml:"s3" json:"s3"`
	GCS        gcsConfig        `yaml:"gcs" json:"gcs"`
	IPFS       ipfsConfig       `yaml:"ipfs" json:"ipfs"`
}

type fileSystemConfig struct {
	DataPath string `yaml:"data_path" json:"data_path" default:"/data" env:"FS_DATA_PATH"`
}

type s3Config struct {
	Hostname        string `yaml:"hostname" json:"hostname" default:"" env:"S3_HOSTNAME"`
	Port            string `yaml:"port" json:"port" default:"" env:"S3_PORT"`
	AccessKeyId     string `yaml:"access_key_id" json:"access_key_id" env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `yaml:"secret_access_key" json:"secret_access_key" env:"S3_SECRET_ACCESS_KEY"`
	Region          string `yaml:"region" json:"region" env:"S3_REGION"`
	Bucket          string `yaml:"bucket" json:"bucket" env:"S3_BUCKET"`
	EnableSSL       bool   `yaml:"enable_ssl" json:"enable_ssl" default:"true" env:"S3_ENABLE_SSL"`
	UsePathStyle    bool   `yaml:"use_path_style" json:"use_path_style" default:"false" env:"S3_ENABLE_PATH_STYLE"`
}

type gcsConfig struct {
	CredentialsFile string `yaml:"credentials_file" json:"credentials_file" default:"" env:"GCS_CREDENTIALS_FILE"`
	Bucket          string `yaml:"bucket" json:"bucket" default:"" env:"GCS_BUCKET"`
}

type ipfsConfig struct {
	Gateway string `yaml:"gateway" json:"gateway" default:"gateway.ipfs.io" env:"IPFS_GATEWAY"`
}

type dbConfig struct {
	Dialect  string `yaml:"-" json:"-" default:"postgres"`
	Host     string `yaml:"host" json:"host" default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" json:"port" default:"5432" env:"DB_PORT"`
	User     string `yaml:"user" json:"user" default:"postgres" env:"DB_USER"`
	Password string `yaml:"password" json:"password" default:"" env:"DB_PASSWORD"`
	Database string `yaml:"db_name" json:"db_name" default:"gilfoyle" env:"DB_NAME"`
}

type redisConfig struct {
	Host     string `yaml:"host" json:"host" default:"localhost" env:"REDIS_HOST"`
	Database string `yaml:"database" json:"database" default:"0" env:"REDIS_DB"`
	Port     string `yaml:"port" json:"port" default:"6379" env:"REDIS_PORT"`
	Password string `yaml:"password" json:"password" default:"" env:"REDIS_PASSWORD"`
}

var config Config

// New creates a new config object
// and load values from environment variables or config file.
// File paths can be both relative and absolute.
func New(files ...string) error {
	err := configor.Load(&config, files...)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig helps you to get configuration data
func GetConfig() *Config {
	return &config
}
