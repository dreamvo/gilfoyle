package mocks

import (
	"github.com/dreamvo/gilfoyle/transcoding"
	"github.com/stretchr/testify/mock"
)

type MockedTranscoder struct {
	mock.Mock
}

func (m *MockedTranscoder) Process() transcoding.IProcess {
	m.Called()
	return new(transcoding.Process)
}

func (m *MockedTranscoder) Run(p transcoding.IProcess) error {
	return m.Called(p).Error(0)
}
