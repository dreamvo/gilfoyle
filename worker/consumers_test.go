package worker

import (
	"context"
	"encoding/json"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/mediafile"
	"github.com/dreamvo/gilfoyle/ent/schema"
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/dreamvo/gilfoyle/x/testutils/mocks"
	_ "github.com/mattn/go-sqlite3"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConsumers(t *testing.T) {
	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() { _ = dbClient.Close() }()

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer func() { _ = os.RemoveAll(gilfoyle.Config.Storage.Filesystem.DataPath) }()

	t.Run("videoTranscodingConsumer", func(t *testing.T) {
		t.Run("should receive one message and succeed", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			params := VideoTranscodingParams{
				MediaUUID: m.ID,
				OriginalFile: transcoding.OriginalFile{
					DurationSeconds: 5.21,
					Format:          "mp4",
					FrameRate:       25,
				},
				VideoWidth:         1280,
				VideoHeight:        720,
				RenditionName:      "360p",
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
				FrameRate:          30,
				TargetBandwidth:    80000,
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)
			transcoderMock := new(mocks.MockedTranscoder)
			storageMock := new(mocks.MockedStorage)

			w := &Worker{
				logger:     loggerMock,
				dbClient:   dbClient,
				storage:    storageMock,
				transcoder: transcoderMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			msgs := make(chan amqp.Delivery)

			storageMock.
				On("Open", filepath.Join(m.ID.String(), transcoding.OriginalFileName)).
				Return(ioutil.NopCloser(strings.NewReader("test")), nil)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
			}).Return()

			transcoderMock.On("Process").Return(new(transcoding.Process))
			transcoderMock.On("Run", mock.Anything).Return(nil)

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			go videoTranscodingConsumer(w, msgs)

			msgs <- delivery

			time.Sleep(200 * time.Millisecond)

			items, err := m.QueryMediaFiles().All(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, 1, len(items))
			assert.Equal(t, "360p", items[0].RenditionName)
			assert.Equal(t, mediafile.MediaType(schema.MediaFileTypeVideo), items[0].MediaType)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
			transcoderMock.AssertExpectations(t)
			storageMock.AssertExpectations(t)
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
			t.Skip()
		})
	})
}
