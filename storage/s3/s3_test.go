package s3_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	assertTest "github.com/stretchr/testify/assert"
)

func TestS3(t *testing.T) {
	assert := assertTest.New(t)

	// fake s3 server
	backend := s3mem.New()
	faker := gofakes3.New(backend)
	ts := httptest.NewServer(faker.Server())
	u, _ := url.Parse(ts.URL)
	host := u.Host
	defer ts.Close()

	gilfoyle.Config.Storage.S3 = config.S3Config{
		Hostname:        host,
		AccessKeyID:     "Q3AM3UQ867SPQQA43P2F",
		SecretAccessKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		Bucket:          "gilfoyle-aws-bucket",
		EnableSSL:       false,
	}

	t.Run("should return error file does not exist", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		_, err = s.Stat(ctx, "doesnotexist")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)
	})

	t.Run("should get metadata of file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		before := time.Now().Add(-1 * time.Second)

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		now := time.Now().Add(2 * time.Second)

		stat, err := s.Stat(ctx, "world")
		assert.NoError(err)

		assert.Equal(int64(5), stat.Size)
		assert.Equal(false, stat.ModifiedTime.Before(before))
		assert.Equal(false, stat.ModifiedTime.After(now))
	})

	t.Run("should create then delete file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		err = s.Delete(ctx, "world")
		assert.NoError(err)

		_, err = s.Stat(ctx, "world")
		assert.EqualError(err, "The specified key does not exist.")
	})

	t.Run("should create then open file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Save(ctx, bytes.NewBufferString("hello"), "world")
		assert.NoError(err)

		f, err := s.Open(ctx, "world")
		assert.NoError(err)
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		assert.NoError(err)
		assert.Equal("hello", string(b))
	})
}
