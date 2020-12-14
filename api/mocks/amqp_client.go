package mocks

import (
	"github.com/dreamvo/gilfoyle/worker"
	"github.com/stretchr/testify/mock"
)

type MockedAMQPClient struct {
	mock.Mock
}

func (m *MockedAMQPClient) Channel() (worker.Channel, error) {
	args := m.Called()
	return nil, args.Error(0)
}

func (m *MockedAMQPClient) Close() error {
	return m.Called().Error(0)
}
