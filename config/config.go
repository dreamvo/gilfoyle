package config

import (
	"github.com/jinzhu/configor"
)

// Config defines the application's settings
type Config struct {
	Services servicesConfig `yaml:"services" json:"services"`
	Settings settingsConfig `yaml:"settings" json:"settings"`
}

type servicesConfig struct {
	IPFS  ipfsConfig  `yaml:"ipfs" json:"ipfs"`
	DB    dbConfig    `yaml:"db" json:"db"`
	Redis redisConfig `yaml:"redis" json:"redis"`
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

type settingsConfig struct {
	ExposeSwaggerUI bool   `yaml:"expose_swagger_ui" json:"expose_swagger_ui" default:"true"`
	MaxFileSize     string `yaml:"max_file_size" json:"max_file_size" default:"50mb"`
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
