package s3_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	assertTest "github.com/stretchr/testify/assert"
)

func TestS3(t *testing.T) {
	assert := assertTest.New(t)
	var err error

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	options := &dockertest.RunOptions{
		Repository: "minio/minio",
		Tag:        "latest",
		Cmd:        []string{"server", "/data"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9000/tcp": []docker.PortBinding{{HostPort: "9000"}},
		},
		Env: []string{"MINIO_ACCESS_KEY=access_key", "MINIO_SECRET_KEY=secret_key"},
	}

	resource, err := pool.RunWithOptions(options, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	endpoint := fmt.Sprintf("localhost:%s", resource.GetPort("9000/tcp"))

	if err := pool.Retry(func() error {
		url := fmt.Sprintf("http://%s/minio/health/live", endpoint)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status code not OK")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	gilfoyle.Config.Storage.S3 = config.S3Config{
		Hostname:        endpoint,
		AccessKeyID:     "access_key",
		SecretAccessKey: "secret_key",
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

	t.Run("should delete the file", func(t *testing.T) {
		s, err := gilfoyle.NewStorage(storage.AmazonS3)
		assert.NoError(err)

		ctx := context.Background()

		err = s.Delete(ctx, "world")
		assert.NoError(err)
	})
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
