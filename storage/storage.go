package storage

import (
	"context"
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage/ipfs"
	"github.com/dreamvo/gilfoyle/storage/s3"
	"github.com/ulule/gostorages"
	"github.com/ulule/gostorages/fs"
	"github.com/ulule/gostorages/gcs"
)

const (
	Filesystem         config.StorageClass = "fs"
	GoogleCloudStorage config.StorageClass = "gcs"
	AmazonS3           config.StorageClass = "s3"
	IPFS               config.StorageClass = "ipfs"
)

// New creates a new storage instance
func New(storageClass config.StorageClass) (gostorages.Storage, error) {
	cfg := config.GetConfig().Storage

	switch storageClass {
	case Filesystem:
		s := fs.NewStorage(fs.Config{Root: cfg.Filesystem.DataPath})
		return s, nil
	case GoogleCloudStorage:
		s, err := gcs.NewStorage(context.Background(), cfg.GCS.CredentialsFile, cfg.GCS.Bucket)
		return s, err
	case AmazonS3:
		return s3.NewStorage(cfg.S3)
	case IPFS:
		return ipfs.NewStorage(cfg.IPFS)
	default:
		return nil, fmt.Errorf("storage class %s does not exist", storageClass)
	}
}
