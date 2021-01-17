package mocks

import (
	"context"
	"github.com/dreamvo/gilfoyle/storage"
	"github.com/stretchr/testify/mock"
	"io"
)

type MockedStorage struct {
	mock.Mock
}

func (m *MockedStorage) Open(_ context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(path)
	return args.Get(0).(io.ReadCloser), args.Error(1)

}

func (m *MockedStorage) Save(_ context.Context, content io.Reader, path string) error {
	args := m.Called(content, path)
	return args.Error(0)

}

func (m *MockedStorage) Stat(_ context.Context, path string) (*storage.Stat, error) {
	args := m.Called(path)
	return args.Get(0).(*storage.Stat), args.Error(1)

}

func (m *MockedStorage) Delete(_ context.Context, path string) error {
	args := m.Called(path)
	return args.Error(0)

}
