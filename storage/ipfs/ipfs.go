package ipfs

import (
	"context"
	"errors"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/ulule/gostorages"
	"io"
)

// Storage is a IPFS storage.
type Storage struct{}

func NewStorage(cfg config.IPFSConfig) (*Storage, error) {
	return nil, errors.New("not implemented yet")
}

// Save saves content to path.
func (s *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	return nil
}

// Stat returns path metadata.
func (s *Storage) Stat(ctx context.Context, path string) (*gostorages.Stat, error) {
	return nil, errors.New("not implemented yet")
}

// Open opens path for reading.
func (s *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented yet")
}

// Delete deletes path.
func (s *Storage) Delete(ctx context.Context, path string) error {
	return errors.New("not implemented yet")
}
