package config

type servicesConfig struct {
	IPFS  ipfsConfig  `yaml:"ipfs"`
	DB    dbConfig    `yaml:"db"`
	Redis redisConfig `yaml:"redis"`
}

type ipfsConfig struct {
	Gateway string `yaml:"gateway" default:"gateway.ipfs.io"`
}

type dbConfig struct {
	Dialect  string `yaml:"dialect"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type redisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

type settingsConfig struct {
	MaxFileSize string `yaml:"maxFileSize" default:"50mb"`
}

// Config defines the application's settings
type Config struct {
	Services servicesConfig `yaml:"services"`
	Settings settingsConfig `yaml:"settings"`
}

// NewConfig a new config object
func NewConfig() *Config {
	c := new(Config)

	// set default values
	c.Services.DB.Dialect = "postgres"
	c.Services.DB.Database = "gilfoyle"
	c.Services.DB.Host = "localhost"
	c.Services.DB.Port = "5432"
	c.Services.DB.User = "postgres"
	c.Services.DB.Password = "secret"

	return c
}

// ParseConfigFile creates a config object from file content
func ParseConfigFile(filepath string) *Config {
	c := NewConfig()

	return c
}
