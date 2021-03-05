package mocks

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type MockedAMQPClient struct {
	mock.Mock
}

func (m *MockedAMQPClient) Channel() (*amqp.Channel, error) {
	args := m.Called()
	return nil, args.Error(1)
}

func (m *MockedAMQPClient) Close() error {
	return m.Called().Error(0)
}

func (m *MockedAMQPClient) IsClosed() bool {
	return m.Called().Bool(0)
}
