package mocks

import "github.com/stretchr/testify/mock"

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
