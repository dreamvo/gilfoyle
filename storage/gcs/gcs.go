package gcs

import (
	"context"
	"github.com/dreamvo/gilfoyle/storage"
	"io"
	"mime"
	"path/filepath"

	gstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Storage is a gcs storage.
type Storage struct {
	bucket *gstorage.BucketHandle
}

// NewStorage returns a new Storage.
func NewStorage(ctx context.Context, credentialsFile, bucket string) (*Storage, error) {
	client, err := gstorage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	return &Storage{bucket: client.Bucket(bucket)}, nil
}

// Save saves content to path.
func (g *Storage) Save(ctx context.Context, content io.Reader, path string) (rerr error) {
	w := g.bucket.Object(path).NewWriter(ctx)
	w.ContentType = mime.TypeByExtension(filepath.Ext(path))

	defer func() {
		if err := w.Close(); err != nil {
			rerr = err
		}
	}()

	if _, err := io.Copy(w, content); err != nil {
		return err
	}

	return rerr
}

// Stat returns path metadata.
func (g *Storage) Stat(ctx context.Context, path string) (*storage.Stat, error) {
	attrs, err := g.bucket.Object(path).Attrs(ctx)
	if err == gstorage.ErrObjectNotExist {
		return nil, storage.ErrNotExist
	} else if err != nil {
		return nil, err
	}

	return &storage.Stat{
		ModifiedTime: attrs.Updated,
		Size:         attrs.Size,
	}, nil
}

// Open opens path for reading.
func (g *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	r, err := g.bucket.Object(path).NewReader(ctx)
	if err == gstorage.ErrObjectNotExist {
		return nil, storage.ErrNotExist
	}
	return r, err
}

// Delete deletes path.
func (g *Storage) Delete(ctx context.Context, path string) error {
	return g.bucket.Object(path).Delete(ctx)
}
