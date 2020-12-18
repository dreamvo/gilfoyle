package s3

import (
	"context"
	"io"
	"mime"
	"path/filepath"

	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/minio/minio-go/v7"
)

type Client interface {
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) (err error)
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
}

// Storage is a s3 storage.
type Storage struct {
	bucket string
	client Client
}

// NewStorage creates a new storage instance
func NewStorage(cfg config.S3Config, client Client) (*Storage, error) {
	ctx := context.Background()

	found, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !found {
		err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil {
			return nil, err
		}
	}

	return &Storage{
		client: client,
		bucket: cfg.Bucket,
	}, err
}

// Save saves content to path.
func (s *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	_, err := s.client.PutObject(ctx, s.bucket, path, content, -1, minio.PutObjectOptions{
		ContentType: mime.TypeByExtension(filepath.Ext(path)),
	})
	return err
}

// Stat returns path metadata.
func (s *Storage) Stat(ctx context.Context, path string) (*storage.Stat, error) {
	stat, err := s.client.StatObject(ctx, s.bucket, path, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	objectStat := &storage.Stat{
		ModifiedTime: stat.LastModified,
		Size:         stat.Size,
	}

	return objectStat, nil
}

// Open opens path for reading.
func (s *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucket, path, minio.GetObjectOptions{})
	return object, err
}

// Delete deletes path.
func (s *Storage) Delete(ctx context.Context, path string) error {
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	return err
}
