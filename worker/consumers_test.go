package worker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dreamvo/gilfoyle/ent/enttest"
	"github.com/dreamvo/gilfoyle/x/testutils/mocks"
	"github.com/floostack/transcoder/ffmpeg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestConsumers(t *testing.T) {
	t.Run("videoTranscodingConsumer", func(t *testing.T) {
		t.Run("should receive one message and succeed", func(t *testing.T) {
			dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			defer func() { _ = dbClient.Close() }()

			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "uuid/test",
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)

			w := &Worker{
				logger:     loggerMock,
				dbClient:   dbClient,
				transcoder: nil,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("MediaUUID", params.MediaUUID.String()),
				zap.String("SourceFilePath", params.SourceFilePath),
			}).Return()

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			go videoTranscodingConsumer(context.Background(), w, delivery)

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

			loggerMock.On("Error", "Unmarshal error", mock.Anything).Return()

			go videoTranscodingConsumer(context.Background(), w, delivery)

			time.Sleep(200 * time.Millisecond)

			loggerMock.AssertExpectations(t)
		})

		t.Run("should fail to send ack", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "uuid/test",
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

			go videoTranscodingConsumer(context.Background(), w, delivery)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("SourceFilePath", params.SourceFilePath),
			}).Return()

			AckMock.On("Ack", mock.Anything, false).Return(errors.New("test"))

			loggerMock.On("Error", "Error trying to send ack", mock.Anything).Return()

			time.Sleep(200 * time.Millisecond)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
		})
	})
}
