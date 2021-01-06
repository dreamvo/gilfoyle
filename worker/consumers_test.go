package worker

import (
	"encoding/json"
	"errors"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/x/testutils/mocks"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

func TestConsumers(t *testing.T) {
	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() { _ = dbClient.Close() }()

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer func() { _ = os.RemoveAll(gilfoyle.Config.Storage.Filesystem.DataPath) }()

	storageDriver, err := gilfoyle.NewStorage(storage.Filesystem)
	if err != nil {
		t.Error(err)
	}

	t.Run("videoTranscodingConsumer", func(t *testing.T) {
		t.Run("should receive one message and succeed", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:          uuid.New(),
				OriginalFilePath:   "uuid/test",
				AudioCodec:         "aac",
				AudioRate:          48000,
				VideoCodec:         "h264",
				Crf:                20,
				KeyframeInterval:   48,
				HlsSegmentDuration: 4,
				HlsPlaylistType:    "vod",
				VideoBitRate:       800000,
				VideoMaxBitRate:    856,
				BufferSize:         1200,
				AudioBitrate:       96000,
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)

			w := &Worker{
				logger:       loggerMock,
				dbClient:     dbClient,
				storage:      storageDriver,
				ffmpegConfig: nil,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			msgs := make(chan amqp.Delivery)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("OriginalFilePath", params.OriginalFilePath),
			}).Return()

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			go videoTranscodingConsumer(w, msgs)

			msgs <- delivery

			time.Sleep(200 * time.Millisecond)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
		})

		t.Run("should fail to unmarshall json", func(t *testing.T) {
			loggerMock := new(mocks.MockedLogger)

			w := &Worker{
				logger: loggerMock,
			}
			delivery := amqp.Delivery{
				Body: []byte(""),
			}

			msgs := make(chan amqp.Delivery)

			loggerMock.On("Error", "Unmarshal error", mock.Anything).Return()

			go videoTranscodingConsumer(w, msgs)

			msgs <- delivery

			time.Sleep(200 * time.Millisecond)

			loggerMock.AssertExpectations(t)
		})

		t.Run("should fail to send ack", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:        uuid.New(),
				OriginalFilePath: "uuid/test",
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)

			w := &Worker{
				logger: loggerMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			msgs := make(chan amqp.Delivery)

			go videoTranscodingConsumer(w, msgs)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("OriginalFilePath", params.OriginalFilePath),
			}).Return()

			AckMock.On("Ack", mock.Anything, false).Return(errors.New("test"))

			loggerMock.On("Error", "Error trying to send ack", mock.Anything).Return()

			msgs <- delivery

			time.Sleep(200 * time.Millisecond)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
		})
	})
}
