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

var (
	Storage storage.Storage
)

// NewStorage creates a new storage instance
func NewStorage(storageClass config.StorageClass) (storage.Storage, error) {
	var err error

	cfg := Config.Storage

	switch storageClass {
	case storage.Filesystem:
		Storage = fs.NewStorage(fs.Config{
			Root: cfg.Filesystem.DataPath,
		})
		return Storage, nil
	case storage.GoogleCloudStorage:
		Storage, err = gcs.NewStorage(context.Background(), cfg.GCS.CredentialsFile, cfg.GCS.Bucket)
		return Storage, err
	case storage.AmazonS3:
		Storage, err = s3.NewStorage(cfg.S3)
		return Storage, err
	//case IPFS:
	//	return ipfs.NewStorage(cfg.IPFS)
	default:
		return nil, fmt.Errorf("storage class %s does not exist", storageClass)
	}
}
