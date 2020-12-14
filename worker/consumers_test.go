package worker

import (
	"encoding/json"
	"errors"
	"github.com/dreamvo/gilfoyle/worker/mocks"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestConsumers(t *testing.T) {
	t.Run("videoTranscodingQueueConsumer", func(t *testing.T) {
		t.Run("should receive one message and succeed", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "uuid/test",
			}

			body, _ := json.Marshal(params)

			loggerMock := new(mocks.MockedLogger)
			AckMock := new(mocks.MockedAcknowledger)

			w := &Worker{
				Logger: loggerMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			msgs := make(chan amqp.Delivery)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("SourceFilePath", params.SourceFilePath),
			}).Return()

			AckMock.On("Ack", mock.Anything, false).Return(nil)

			go videoTranscodingQueueConsumer(w, msgs)

			msgs <- delivery

			time.Sleep(200 * time.Millisecond)

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
		})

		t.Run("should fail to unmarshall json", func(t *testing.T) {
			loggerMock := new(mocks.MockedLogger)

			w := &Worker{
				Logger: loggerMock,
			}
			delivery := amqp.Delivery{
				Body: []byte(""),
			}

			msgs := make(chan amqp.Delivery)

			loggerMock.On("Error", "Unmarshal error", mock.Anything).Return()

			go videoTranscodingQueueConsumer(w, msgs)

			msgs <- delivery

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
				Logger: loggerMock,
			}
			delivery := amqp.Delivery{
				Body:         body,
				Acknowledger: AckMock,
			}

			msgs := make(chan amqp.Delivery)

			go videoTranscodingQueueConsumer(w, msgs)

			loggerMock.On("Info", "Received a message", []zap.Field{
				zap.String("SourceFilePath", params.SourceFilePath),
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
