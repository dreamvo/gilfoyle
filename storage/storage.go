package storage

import (
	"context"
	"errors"
	"github.com/dreamvo/gilfoyle/config"
	"io"
	"time"
)

const (
	Filesystem         config.StorageDriver = "fs"
	GoogleCloudStorage config.StorageDriver = "gcs"
	AmazonS3           config.StorageDriver = "s3"
	_                  config.StorageDriver = "openstack"
)

// Storage is the storage interface.
type Storage interface {
	Save(ctx context.Context, content io.Reader, path string) error
	Stat(ctx context.Context, path string) (*Stat, error)
	Open(ctx context.Context, path string) (io.ReadCloser, error)
	Delete(ctx context.Context, path string) error
}

// Stat contains metadata about content stored in storage.
type Stat struct {
	ModifiedTime time.Time
	Size         int64
}

// ErrNotExist is a sentinel error returned by the Open and the Stat methods.
var ErrNotExist = errors.New("file does not exist")
