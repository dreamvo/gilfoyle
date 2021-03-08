package worker

import (
	"encoding/json"
	"errors"
	"github.com/dreamvo/gilfoyle/x/testutils/mocks"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducers(t *testing.T) {
	t.Run("EncodingEntrypointProducer", func(t *testing.T) {
		t.Run("should publish a new message", func(t *testing.T) {
			params := EncodingEntrypointParams{
				MediaUUID: uuid.New(),
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", EncodingEntrypointQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(nil)

			err = EncodingEntrypointProducer(ch, params)
			assert.NoError(t, err)

			ch.AssertExpectations(t)
		})

		t.Run("should publish a new message with AMQP error", func(t *testing.T) {
			params := EncodingEntrypointParams{
				MediaUUID: uuid.New(),
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", EncodingEntrypointQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(errors.New("test"))

			err = EncodingEntrypointProducer(ch, params)
			assert.EqualError(t, err, "test")

			ch.AssertExpectations(t)
		})
	})

	t.Run("HlsVideoEncodingProducer", func(t *testing.T) {
		t.Run("should publish a new message", func(t *testing.T) {
			params := HlsVideoEncodingParams{
				MediaFileUUID:      uuid.New(),
				KeyframeInterval:   48,
				HlsPlaylistType:    "vod",
				HlsSegmentDuration: 4,
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", HlsVideoEncodingQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(nil)

			err = HlsVideoEncodingProducer(ch, params)
			assert.NoError(t, err)

			ch.AssertExpectations(t)
		})

		t.Run("should publish a new message with AMQP error", func(t *testing.T) {
			params := HlsVideoEncodingParams{
				MediaFileUUID:      uuid.New(),
				KeyframeInterval:   48,
				HlsPlaylistType:    "vod",
				HlsSegmentDuration: 4,
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", HlsVideoEncodingQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(errors.New("test"))

			err = HlsVideoEncodingProducer(ch, params)
			assert.EqualError(t, err, "test")

			ch.AssertExpectations(t)
		})
	})

	t.Run("EncodingFinalizerProducer", func(t *testing.T) {
		t.Run("should publish a new message", func(t *testing.T) {
			params := EncodingFinalizerParams{
				MediaUUID: uuid.New(),
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", EncodingFinalizerQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(nil)

			err = EncodingFinalizerProducer(ch, params)
			assert.NoError(t, err)

			ch.AssertExpectations(t)
		})

		t.Run("should publish a new message with AMQP error", func(t *testing.T) {
			params := EncodingFinalizerParams{
				MediaUUID: uuid.New(),
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", EncodingFinalizerQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(errors.New("test"))

			err = EncodingFinalizerProducer(ch, params)
			assert.EqualError(t, err, "test")

			ch.AssertExpectations(t)
		})
	})
}
