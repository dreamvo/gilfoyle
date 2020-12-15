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
}

func (m *MockedLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockedLogger) Fatal(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockedLogger) With(fields ...zap.Field) *zap.Logger {
	m.Called(fields)
	return new(zap.Logger)
}
