package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
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
