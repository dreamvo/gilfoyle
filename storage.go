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
func NewStorage(storageClass config.StorageClass) (storage.Storage, error) {
	cfg := Config.Storage

	switch storageClass {
	case storage.Filesystem:
		return fs.NewStorage(fs.Config{
			Root: cfg.Filesystem.DataPath,
		}), nil
	case storage.GoogleCloudStorage:
		return gcs.NewStorage(context.Background(), cfg.GCS.CredentialsFile, cfg.GCS.Bucket)
	case storage.AmazonS3:
		return s3.NewStorage(cfg.S3)
	//case IPFS:
	//	return ipfs.NewStorage(cfg.IPFS)
	default:
		return nil, fmt.Errorf("storage class %s does not exist", storageClass)
	}
}
