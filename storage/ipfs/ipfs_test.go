package ipfs

import (
	"bytes"
	"context"
	"github.com/dreamvo/gilfoyle/config"
	assertTest "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestIPFS(t *testing.T) {
	// TODO(sundowndev): test this package
	t.Skip()

	assert := assertTest.New(t)

	cfg := config.IPFSConfig{
		Gateway: "gateway.ipfs.io",
	}

	t.Run("should return error file does not exist", func(t *testing.T) {
		storage, err := NewStorage(cfg)
		assert.NoError(err)

		ctx := context.Background()

		_, err = storage.Stat(ctx, "doesnotexist")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create file", func(t *testing.T) {
		storage, err := NewStorage(cfg)
		assert.NoError(err)

		ctx := context.Background()

		err = storage.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)
	})

	t.Run("should get metadata of file", func(t *testing.T) {
		storage, err := NewStorage(cfg)
		assert.NoError(err)

		ctx := context.Background()

		before := time.Now()

		err = storage.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		now := time.Now().Add(5 * time.Second)

		stat, err := storage.Stat(ctx, "world")
		assert.NoError(err)

		assert.Equal(int64(5), stat.Size)
		assert.Equal(false, stat.ModifiedTime.Before(before))
		assert.Equal(false, stat.ModifiedTime.After(now))
	})

	t.Run("should create then delete file", func(t *testing.T) {
		storage, err := NewStorage(cfg)
		assert.NoError(err)

		ctx := context.Background()

		err = storage.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		err = storage.Delete(ctx, "world")
		assert.NoError(err)

		_, err = storage.Stat(ctx, "world")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create then open file", func(t *testing.T) {
		storage, err := NewStorage(cfg)
		assert.NoError(err)

		ctx := context.Background()

		err = storage.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		f, err := storage.Open(ctx, "world")
		assert.NoError(err)
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		assert.NoError(err)
		assert.Equal("hello", string(b))
	})
}
