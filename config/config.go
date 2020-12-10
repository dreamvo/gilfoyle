package config

// StorageClass is a kind of storage backend
type StorageClass string

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
	ExposeSwaggerUI bool           `yaml:"expose_swagger_ui" json:"expose_swagger_ui" default:"true"`
	MaxFileSize     int64          `yaml:"max_file_size" json:"max_file_size" default:"524288000"`
	Debug           bool           `yaml:"debug" json:"debug" default:"false" env:"APP_DEBUG"`
	Worker          WorkerSettings `yaml:"worker" json:"worker"`
}

type StorageConfig struct {
	Class      string           `yaml:"class" json:"class" default:"fs" env:"STORAGE_CLASS"`
	Filesystem FileSystemConfig `yaml:"fs" json:"fs"`
	S3         S3Config         `yaml:"s3" json:"s3"`
	GCS        GCSConfig        `yaml:"gcs" json:"gcs"`
	IPFS       IPFSConfig       `yaml:"ipfs" json:"ipfs"`
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

type IPFSConfig struct {
	Gateway string `yaml:"gateway" json:"gateway" default:"gateway.ipfs.io" env:"IPFS_GATEWAY"`
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
	Concurrency uint `yaml:"concurrency" json:"concurrency" default:"3"`
}
