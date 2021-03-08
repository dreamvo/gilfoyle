package worker

import (
	"context"
	"encoding/json"
	"github.com/dreamvo/gilfoyle"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/ent/media"
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
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestConsumers(t *testing.T) {
	dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer func() { _ = dbClient.Close() }()

	gilfoyle.Config.Storage.Filesystem.DataPath = "./data"
	defer func() { _ = os.RemoveAll(gilfoyle.Config.Storage.Filesystem.DataPath) }()

	t.Run("hlsVideoEncodingConsumer", func(t *testing.T) {
		t.Run("should receive one message and succeed", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			mediaFile, err := dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("test-rendition").
				SetMediaType(schema.MediaFileTypeVideo).
				SetFormat("hls").
				SetStatus(mediafile.StatusProcessing).
				SetEntryFile(transcoding.HLSPlaylistFilename).
				SetMimetype(transcoding.HLSPlaylistMimeType).
				SetTargetBandwidth(928000).
				SetVideoBitrate(800000).
				SetAudioBitrate(128000).
				SetVideoCodec("h264").
				SetAudioCodec("aac").
				SetResolutionWidth(640).
				SetResolutionHeight(360).
				SetDurationSeconds(5).
				SetFramerate(25).
				Save(context.Background())
			assert.NoError(t, err)

			params := HlsVideoEncodingParams{
				MediaFileUUID:      mediaFile.ID,
				KeyframeInterval:   48,
				HlsPlaylistType:    "vod",
				HlsSegmentDuration: 4,
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

			loggerMock.On("Info", "Received HLS video encoding message", []zap.Field{
				zap.String("MediaFileUUID", params.MediaFileUUID.String()),
			}).Return()

			storageMock.
				On("Open", filepath.Join(m.ID.String(), transcoding.OriginalFileName)).
				Return(ioutil.NopCloser(strings.NewReader("test")), nil)

			transcoderMock.On("Process").Return()
			transcoderMock.On("Run", mock.Anything).Return(nil)

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			hlsVideoEncodingConsumer(w, delivery)

			items, err := m.QueryMediaFiles().All(context.Background())
			assert.NoError(t, err)

			assert.Equal(t, 1, len(items))
			assert.Equal(t, mediafile.StatusReady, items[0].Status)
			assert.Equal(t, "Encoding job succeeded", items[0].Message)

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

			loggerMock.On("Error", "Unmarshal error", mock.Anything).Return()

			hlsVideoEncodingConsumer(w, delivery)

			loggerMock.AssertExpectations(t)
		})
	})

	t.Run("encodingFinalizerConsumer", func(t *testing.T) {
		t.Run("should finalize a media with no rendition", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			params := EncodingFinalizerParams{
				MediaUUID: m.ID,
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)
			storageMock := new(mocks.MockedStorage)

			w := &Worker{
				logger:   loggerMock,
				dbClient: dbClient,
				storage:  storageMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			loggerMock.On("Info", "Received encoding finalizer message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
			}).Return()

			storageMock.
				On(
					"Save",
					strings.NewReader(`#EXTM3U
#EXT-X-VERSION:3
`),
					path.Join(m.ID.String(), transcoding.HLSMasterPlaylistFilename),
				).
				Return(nil)

			AckMock.On("Ack", uint64(0), false).Return(nil)

			encodingFinalizerConsumer(w, delivery)

			m, err = dbClient.Media.Get(context.Background(), m.ID)
			assert.NoError(t, err)
			assert.Equal(t, media.StatusErrored, m.Status)
			assert.Equal(t, "Media doesn't have any rendition", m.Message)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
			storageMock.AssertExpectations(t)
		})

		t.Run("should finalize a media with ready rendition", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("test-rendition").
				SetMediaType(schema.MediaFileTypeVideo).
				SetFormat("hls").
				SetStatus(mediafile.StatusReady).
				SetEntryFile(transcoding.HLSPlaylistFilename).
				SetMimetype(transcoding.HLSPlaylistMimeType).
				SetTargetBandwidth(928000).
				SetVideoBitrate(800000).
				SetAudioBitrate(128000).
				SetVideoCodec("h264").
				SetAudioCodec("aac").
				SetResolutionWidth(640).
				SetResolutionHeight(360).
				SetDurationSeconds(5).
				SetFramerate(25).
				Save(context.Background())
			assert.NoError(t, err)

			params := EncodingFinalizerParams{
				MediaUUID: m.ID,
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)
			storageMock := new(mocks.MockedStorage)

			w := &Worker{
				logger:   loggerMock,
				dbClient: dbClient,
				storage:  storageMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			loggerMock.On("Info", "Received encoding finalizer message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
			}).Return()

			storageMock.
				On(
					"Save",
					strings.NewReader(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=928000,RESOLUTION=640x360,NAME="test-rendition",FRAME-RATE=25.000
test-rendition/index.m3u8
`),
					path.Join(m.ID.String(), transcoding.HLSMasterPlaylistFilename),
				).
				Return(nil)

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			encodingFinalizerConsumer(w, delivery)

			m, err = dbClient.Media.Get(context.Background(), m.ID)
			assert.NoError(t, err)
			assert.Equal(t, media.StatusReady, m.Status)
			assert.Equal(t, "One or more rendition succeeded. Media is available for streaming", m.Message)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
			storageMock.AssertExpectations(t)
		})

		t.Run("should finalize a media with a errored rendition", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("test-rendition").
				SetMediaType(schema.MediaFileTypeVideo).
				SetFormat("hls").
				SetStatus(mediafile.StatusErrored).
				SetEntryFile(transcoding.HLSPlaylistFilename).
				SetMimetype(transcoding.HLSPlaylistMimeType).
				SetTargetBandwidth(928000).
				SetVideoBitrate(800000).
				SetAudioBitrate(128000).
				SetVideoCodec("h264").
				SetAudioCodec("aac").
				SetResolutionWidth(640).
				SetResolutionHeight(360).
				SetDurationSeconds(5).
				SetFramerate(25).
				Save(context.Background())
			assert.NoError(t, err)

			params := EncodingFinalizerParams{
				MediaUUID: m.ID,
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)
			storageMock := new(mocks.MockedStorage)

			w := &Worker{
				logger:   loggerMock,
				dbClient: dbClient,
				storage:  storageMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			loggerMock.On("Info", "Received encoding finalizer message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
			}).Return()

			storageMock.
				On(
					"Save",
					strings.NewReader(`#EXTM3U
#EXT-X-VERSION:3
`),
					path.Join(m.ID.String(), transcoding.HLSMasterPlaylistFilename),
				).
				Return(nil)

			AckMock.On("Ack", uint64(0), false).Return(nil)

			encodingFinalizerConsumer(w, delivery)

			m, err = dbClient.Media.Get(context.Background(), m.ID)
			assert.NoError(t, err)
			assert.Equal(t, media.StatusErrored, m.Status)
			assert.Equal(t, "All encoding jobs failed", m.Message)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
			storageMock.AssertExpectations(t)
		})

		t.Run("should requeue if a media is in processing state", func(t *testing.T) {
			m, err := dbClient.Media.
				Create().
				SetTitle("my video").
				SetStatus(schema.MediaStatusProcessing).
				SetOriginalFilename(transcoding.OriginalFileName).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("test-rendition").
				SetMediaType(schema.MediaFileTypeVideo).
				SetFormat("hls").
				SetStatus(mediafile.StatusProcessing).
				SetEntryFile(transcoding.HLSPlaylistFilename).
				SetMimetype(transcoding.HLSPlaylistMimeType).
				SetTargetBandwidth(928000).
				SetVideoBitrate(800000).
				SetAudioBitrate(128000).
				SetVideoCodec("h264").
				SetAudioCodec("aac").
				SetResolutionWidth(640).
				SetResolutionHeight(360).
				SetDurationSeconds(5).
				SetFramerate(25).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("test-rendition").
				SetMediaType(schema.MediaFileTypeVideo).
				SetFormat("hls").
				SetStatus(mediafile.StatusErrored).
				SetEntryFile(transcoding.HLSPlaylistFilename).
				SetMimetype(transcoding.HLSPlaylistMimeType).
				SetTargetBandwidth(928000).
				SetVideoBitrate(800000).
				SetAudioBitrate(128000).
				SetVideoCodec("h264").
				SetAudioCodec("aac").
				SetResolutionWidth(640).
				SetResolutionHeight(360).
				SetDurationSeconds(5).
				SetFramerate(25).
				Save(context.Background())
			assert.NoError(t, err)

			_, err = dbClient.MediaFile.
				Create().
				SetMedia(m).
				SetRenditionName("test-rendition").
				SetMediaType(schema.MediaFileTypeVideo).
				SetFormat("hls").
				SetStatus(mediafile.StatusReady).
				SetEntryFile(transcoding.HLSPlaylistFilename).
				SetMimetype(transcoding.HLSPlaylistMimeType).
				SetTargetBandwidth(928000).
				SetVideoBitrate(800000).
				SetAudioBitrate(128000).
				SetVideoCodec("h264").
				SetAudioCodec("aac").
				SetResolutionWidth(640).
				SetResolutionHeight(360).
				SetDurationSeconds(5).
				SetFramerate(25).
				Save(context.Background())
			assert.NoError(t, err)

			params := EncodingFinalizerParams{
				MediaUUID: m.ID,
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)
			storageMock := new(mocks.MockedStorage)

			w := &Worker{
				logger:   loggerMock,
				dbClient: dbClient,
				storage:  storageMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			loggerMock.On("Info", "Received encoding finalizer message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
			}).Return()

			storageMock.
				On(
					"Save",
					strings.NewReader(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=928000,RESOLUTION=640x360,NAME="test-rendition",FRAME-RATE=25.000
test-rendition/index.m3u8
`),
					path.Join(m.ID.String(), transcoding.HLSMasterPlaylistFilename),
				).
				Return(nil)

			AckMock.On("Ack", uint64(0), false).Return(nil)

			encodingFinalizerConsumer(w, delivery)

			m, err = dbClient.Media.Get(context.Background(), m.ID)
			assert.NoError(t, err)
			assert.Equal(t, media.StatusProcessing, m.Status)
			assert.Equal(t, "Media is not yet available for streaming", m.Message)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
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

			loggerMock.On("Error", "Unmarshal error", mock.Anything).Return()

			encodingFinalizerConsumer(w, delivery)

			loggerMock.AssertExpectations(t)
		})
	})
}
