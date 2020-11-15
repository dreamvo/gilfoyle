package gilfoyle

import (
	"context"
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage/s3"
	"github.com/ulule/gostorages"
	"github.com/ulule/gostorages/fs"
	"github.com/ulule/gostorages/gcs"
)

const (
	Filesystem         config.StorageClass = "fs"
	GoogleCloudStorage config.StorageClass = "gcs"
	AmazonS3           config.StorageClass = "s3"
	_                  config.StorageClass = "ipfs"
	_                  config.StorageClass = "openstack"
)

var (
	Storage gostorages.Storage
)

// NewStorage creates a new storage instance
func NewStorage(storageClass config.StorageClass) (gostorages.Storage, error) {
	var err error

	cfg := Config.Storage

	switch storageClass {
	case Filesystem:
		Storage = fs.NewStorage(fs.Config{Root: cfg.Filesystem.DataPath})
		return Storage, nil
	case GoogleCloudStorage:
		Storage, err = gcs.NewStorage(context.Background(), cfg.GCS.CredentialsFile, cfg.GCS.Bucket)
		return Storage, err
	case AmazonS3:
		Storage, err = s3.NewStorage(cfg.S3)
		return Storage, err
	//case IPFS:
	//	return ipfs.NewStorage(cfg.IPFS)
	default:
		return nil, fmt.Errorf("storage class %s does not exist", storageClass)
	}
}
