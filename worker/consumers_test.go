package worker

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

type MockedLogger struct {
	mock.Mock
}

func (m *MockedLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
	return
}

func (m *MockedLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
	return
}

func (m *MockedLogger) Fatal(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
	return
}

type MockedAcknowledger struct {
	mock.Mock
}

func (m *MockedAcknowledger) Ack(tag uint64, multiple bool) error {
	args := m.Called(tag, multiple)
	return args.Error(0)
}

func (m *MockedAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	args := m.Called(tag, multiple, requeue)
	return args.Error(0)
}

func (m *MockedAcknowledger) Reject(tag uint64, requeue bool) error {
	args := m.Called(tag, requeue)
	return args.Error(0)
}

func TestConsumers(t *testing.T) {
	t.Run("videoTranscodingQueueConsumer", func(t *testing.T) {
		t.Run("should receive one message and succeed", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "uuid/test",
			}

			body, _ := json.Marshal(params)

			loggerMock := new(MockedLogger)
			AckMock := new(MockedAcknowledger)

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

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
		})

		t.Run("should fail to unmarshall json", func(t *testing.T) {
			loggerMock := new(MockedLogger)

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

			loggerMock.AssertExpectations(t)
		})

		t.Run("should fail to send ack", func(t *testing.T) {
			params := VideoTranscodingParams{
				MediaUUID:      uuid.New(),
				SourceFilePath: "uuid/test",
			}

			body, _ := json.Marshal(params)

			loggerMock := new(MockedLogger)
			AckMock := new(MockedAcknowledger)

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

			loggerMock.AssertExpectations(t)
			AckMock.AssertExpectations(t)
		})
	})
}
