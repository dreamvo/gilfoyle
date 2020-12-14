package worker_test

import (
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/dreamvo/gilfoyle/x/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorker(t *testing.T) {
	container := testutils.CreateRabbitMQContainer(t, "guest", "guest")
	defer testutils.StopContainer(t, container)

	opts := worker.Options{
		Host:        container.Host,
		Port:        container.DefaultPort(),
		Username:    "guest",
		Password:    "guest",
		Concurrency: 1,
	}

	t.Run("should create new client & declare queues", func(t *testing.T) {
		w, err := worker.New(opts)
		assert.NoError(t, err)
		defer testutils.CloseWorker(t, w)

		err = w.Init()
		assert.NoError(t, err)
	})

	t.Run("should fail to close connection", func(t *testing.T) {
		w, err := worker.New(opts)
		assert.NoError(t, err)

		assert.False(t, w.Client.IsClosed())

		err = w.Close()
		assert.NoError(t, err)

		assert.True(t, w.Client.IsClosed())

		err = w.Close()
		assert.EqualError(t, err, "Exception (504) Reason: \"channel/connection is not open\"")
	})

	t.Run("should fail to connect", func(t *testing.T) {
		_, err := worker.New(worker.Options{
			Host: "127.0.0.1",
			Port: 1000,
		})
		assert.EqualError(t, err, "dial tcp 127.0.0.1:1000: connect: connection refused")
	})

	t.Run("should start consuming queues", func(t *testing.T) {
		w, err := worker.New(opts)
		assert.NoError(t, err)
		defer testutils.CloseWorker(t, w)

		err = w.Init()
		assert.NoError(t, err)

		err = w.Consume()
		assert.NoError(t, err)

		ch, err := w.Client.Channel()
		assert.NoError(t, err)

		q, err := ch.QueueInspect(worker.VideoTranscodingQueue)
		assert.NoError(t, err)

		assert.Equal(t, 0, q.Messages)
		assert.Equal(t, 1, q.Consumers)
	})

	t.Run("should fail to declare queues", func(t *testing.T) {
		w, err := worker.New(opts)
		assert.NoError(t, err)

		err = w.Close()
		assert.NoError(t, err)

		err = w.Init()
		assert.EqualError(t, err, "Exception (504) Reason: \"channel/connection is not open\"")
	})

	t.Run("should fail to start consuming queues", func(t *testing.T) {
		w, err := worker.New(opts)
		assert.NoError(t, err)

		err = w.Init()
		assert.NoError(t, err)

		err = w.Close()
		assert.NoError(t, err)

		err = w.Consume()
		assert.EqualError(t, err, "error creating channel: &{%!e(int=504) %!e(string=channel/connection is not open) %!e(bool=false) %!e(bool=false)}")
	})
}
