package config

import (
	"github.com/jinzhu/configor"
)

// Config defines the application's settings
type Config struct {
	Services servicesConfig `yaml:"services"`
	Settings settingsConfig `yaml:"settings"`
}

type servicesConfig struct {
	IPFS  ipfsConfig  `yaml:"ipfs"`
	DB    dbConfig    `yaml:"db"`
	Redis redisConfig `yaml:"redis"`
}

type ipfsConfig struct {
	Gateway string `yaml:"gateway" default:"gateway.ipfs.io" env:"IPFS_GATEWAY"`
}

type dbConfig struct {
	Dialect  string `yaml:"dialect" default:"postgres" env:"DB_HOST"`
	Host     string `yaml:"host" default:"localhost" env:"DB_HOST"`
	Port     string `yaml:"port" default:"5432" env:"DB_PORT"`
	User     string `yaml:"user" default:"postgres" env:"DB_USER"`
	Password string `yaml:"password" default:"secret" env:"DB_PASSWORD"`
	Database string `yaml:"database" default:"gilfoyle" env:"DB_NAME"`
}

type redisConfig struct {
	Host     string `yaml:"host" default:"localhost" env:"REDIS_HOST"`
	Port     string `yaml:"port" default:"6379" env:"REDIS_PORT"`
	Password string `yaml:"password" default:"" env:"REDIS_PASSWORD"`
}

type settingsConfig struct {
	MaxFileSize string `yaml:"maxFileSize" default:"50mb"`
}

var c Config

// NewConfig a new config object
func NewConfig() *Config {
	_ = configor.Load(&c)

	return &c
}

// GetConfig helps you to get configuration data
func GetConfig() *Config {
	return &c
}
