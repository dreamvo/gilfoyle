package worker

import (
	"encoding/json"
	"errors"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/x/testutils/mocks"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducers(t *testing.T) {
	t.Run("HlsVideoEncodingProducer", func(t *testing.T) {
		t.Run("should publish a new message", func(t *testing.T) {
			params := HlsVideoEncodingParams{
				MediaUUID: uuid.New(),
				OriginalFile: transcoding.OriginalFile{
					Filepath:        "uuid/test",
					DurationSeconds: 5.21,
					Format:          "mp4",
					FrameRate:       25,
				},
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
				MediaUUID: uuid.New(),
				OriginalFile: transcoding.OriginalFile{
					Filepath:        "uuid/test",
					DurationSeconds: 5.21,
					Format:          "mp4",
					FrameRate:       25,
				},
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

	t.Run("MediaEncodingCallbackProducer", func(t *testing.T) {
		t.Run("should publish a new message", func(t *testing.T) {
			params := MediaEncodingCallbackParams{
				MediaUUID:       uuid.New(),
				MediaFilesCount: 1,
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", MediaEncodingCallbackQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(nil)

			err = MediaEncodingCallbackProducer(ch, params)
			assert.NoError(t, err)

			ch.AssertExpectations(t)
		})

		t.Run("should publish a new message with AMQP error", func(t *testing.T) {
			params := MediaEncodingCallbackParams{
				MediaUUID:       uuid.New(),
				MediaFilesCount: 1,
			}

			body, err := json.Marshal(params)
			assert.NoError(t, err)

			ch := new(mocks.MockedChannel)

			ch.On("Publish", "", MediaEncodingCallbackQueue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         body,
			}).Return(errors.New("test"))

			err = MediaEncodingCallbackProducer(ch, params)
			assert.EqualError(t, err, "test")

			ch.AssertExpectations(t)
		})
	})
}
