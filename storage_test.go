package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("should use non-existing storage class", func(t *testing.T) {
		_, err := NewStorage("test")
		assert.EqualError(t, err, "storage class test does not exist")
	})

	t.Run("should initialize Filesystem storage", func(t *testing.T) {
		_, err := NewStorage(storage.Filesystem)
		assert.NoError(t, err)
	})

	t.Run("should initialize S3 storage", func(t *testing.T) {
		Config.Storage.S3 = config.S3Config{
			Hostname:        "play.min.io",
			AccessKeyID:     "Q3AM3UQ867SPQQA43P2F",
			SecretAccessKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
			Bucket:          uuid.New().String(),
			EnableSSL:       true,
		}

		_, err := NewStorage(storage.AmazonS3)
		assert.NoError(t, err)
	})
}
