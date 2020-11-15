package fs

import (
	"context"
	"github.com/dreamvo/gilfoyle/storage"
	"io"
	"os"
	"path/filepath"
)

// Storage is a filesystem storage.
type Storage struct {
	root string
}

// NewStorage returns a new filesystem storage.
func NewStorage(cfg Config) *Storage {
	return &Storage{root: cfg.Root}
}

// Config is the configuration for Storage.
type Config struct {
	Root string
}

func (fs *Storage) abs(path string) string {
	return filepath.Join(fs.root, path)
}

// Save saves content to path.
func (fs *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	abs := fs.abs(path)
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return err
	}

	w, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := io.Copy(w, content); err != nil {
		return err
	}
	return nil
}

// Stat returns path metadata.
func (fs *Storage) Stat(ctx context.Context, path string) (*storage.Stat, error) {
	fi, err := os.Stat(fs.abs(path))
	if os.IsNotExist(err) {
		return nil, storage.ErrNotExist
	} else if err != nil {
		return nil, err
	}

	return &storage.Stat{
		ModifiedTime: fi.ModTime(),
		Size:         fi.Size(),
	}, nil
}

// Open opens path for reading.
func (fs *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	f, err := os.Open(fs.abs(path))
	if os.IsNotExist(err) {
		return nil, storage.ErrNotExist
	}
	return f, err
}

// Delete deletes path.
func (fs *Storage) Delete(ctx context.Context, path string) error {
	return os.Remove(fs.abs(path))
}
