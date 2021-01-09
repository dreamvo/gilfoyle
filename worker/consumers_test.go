package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/x/testutils/mocks"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
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
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			originalPath := fmt.Sprintf("%s/%s/%s", gilfoyle.Config.Storage.Filesystem.DataPath, m.ID.String(), transcoding.OriginalFileName)

			f, err := os.Open("../x/testutils/fixtures/SampleVideo_1280x720_1mb.mp4")
			assert.NoError(t, err)

			err = storageDriver.Save(context.Background(), f, originalPath)
			assert.NoError(t, err)

			params := VideoTranscodingParams{
				MediaUUID: m.ID,
				OriginalFile: transcoding.OriginalFile{
					DurationSeconds: 5.21,
					Format:          "mp4",
					FrameRate:       25,
				},
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
			transcoderMock := new(mocks.MockedTranscoder)

			w := &Worker{
				logger:     loggerMock,
				dbClient:   dbClient,
				storage:    storageDriver,
				transcoder: transcoderMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			msgs := make(chan amqp.Delivery)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
			}).Return()

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			transcoderMock.On("Process").Return()
			transcoderMock.On("Run", transcoding.Process{}).Return(nil)

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
				MediaUUID: uuid.New(),
				OriginalFile: transcoding.OriginalFile{
					Filepath:        "uuid/test",
					DurationSeconds: 5.21,
					Format:          "mp4",
					FrameRate:       25,
				},
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
				zap.String("OriginalFilePath", params.OriginalFile.Filepath),
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
