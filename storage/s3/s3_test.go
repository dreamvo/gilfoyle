package s3_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/storage/s3"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	assertTest "github.com/stretchr/testify/assert"
)

type ClientImpl struct {
	s3.Client
}

func (c *ClientImpl) BucketExistsMock(ctx context.Context, bucketName string) (bool, error) {
	return true, errors.New("Error")
}

func (c *ClientImpl) MakeBucketMock(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) (err error) {
	return errors.New("Error")
}

func (c *ClientImpl) PutObjectMock(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
	opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	return minio.UploadInfo{
		Bucket:           "string",
		Key:              "string",
		ETag:             "string",
		Size:             64,
		LastModified:     time.Now(),
		Location:         "string",
		VersionID:        "string",
		Expiration:       time.Now(),
		ExpirationRuleID: "string",
	}, errors.New("Error")
}

func (c *ClientImpl) StatObjectMock(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	return minio.ObjectInfo{}, errors.New("Error")
}

func (c *ClientImpl) GetObjectMock(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return &minio.Object{}, errors.New("Error")
}

func (c *ClientImpl) RemoveObjectMock(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) error {
	return nil
}

func TestS3(t *testing.T) {
	assert := assertTest.New(t)

	client := &ClientImpl{}

	cfg := config.S3Config{
		Hostname:        "",
		AccessKeyID:     "",
		SecretAccessKey: "",
		Bucket:          uuid.New().String(),
		EnableSSL:       true,
	}

	gilfoyle.Config.Storage.S3 = config.S3Config{
		Hostname:        "play.min.io",
		AccessKeyID:     "Q3AM3UQ867SPQQA43P2F",
		SecretAccessKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		Bucket:          uuid.New().String(),
		EnableSSL:       true,
	}

	t.Run("should return error file does not exist", func(t *testing.T) {
		s, err := s3.NewStorage(cfg, client)
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
