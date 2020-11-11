package storage

import (
	"context"
	"errors"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/ulule/gostorages"
	"github.com/ulule/gostorages/fs"
	"github.com/ulule/gostorages/gcs"
	"github.com/ulule/gostorages/s3"
)

const (
	Filesystem         config.StorageClass = "fs"
	GoogleCloudStorage config.StorageClass = "gcs"
	AmazonS3           config.StorageClass = "s3"
	IPFS               config.StorageClass = "ipfs"
)

// New creates a new storage instance
func New(k config.StorageClass) (gostorages.Storage, error) {
	switch k {
	case Filesystem:
		s := fs.NewStorage(fs.Config{Root: config.GetConfig().Storage.CachePath})
		return s, nil
	case GoogleCloudStorage:
		s, err := gcs.NewStorage(context.Background(), config.GetConfig().Storage.GCS.CredentialsFile, config.GetConfig().Storage.GCS.Bucket)
		return s, err
	case AmazonS3:
		s, err := s3.NewStorage(s3.Config{
			AccessKeyID:     config.GetConfig().Storage.S3.AccessKeyId,
			SecretAccessKey: config.GetConfig().Storage.S3.SecretAccessKey,
			Region:          config.GetConfig().Storage.S3.Region,
			Bucket:          config.GetConfig().Storage.S3.Bucket,
		})
		return s, err
	case IPFS:
		return nil, errors.New("not implemented yet")
	default:
		return nil, errors.New("wrong storage kind")
	}
}
