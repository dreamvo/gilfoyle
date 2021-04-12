package gilfoyle

import (
	"context"
	"fmt"

	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/storage/fs"
	"github.com/dreamvo/gilfoyle/storage/gcs"
	"github.com/dreamvo/gilfoyle/storage/s3"
)

// NewStorage creates a new storage instance
func NewStorage(cfg config.Config) (storage.Storage, error) {
	driver := config.StorageDriver(cfg.Storage.Driver)

	switch driver {
	case storage.Filesystem:
		return fs.NewStorage(fs.Config{
			Root: cfg.Storage.Filesystem.DataPath,
		}), nil
	case storage.GoogleCloudStorage:
		return gcs.NewStorage(context.Background(), cfg.Storage.GCS.CredentialsFile, cfg.Storage.GCS.Bucket)
	case storage.AmazonS3:
		return s3.NewStorage(cfg.Storage.S3)
	default:
		return nil, fmt.Errorf("storage driver %s does not exist", driver)
	}
}
