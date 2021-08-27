package storage

import (
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (sm *StorageMock) Connect(connectionUrl string) error {
	args := sm.Called(connectionUrl)
	return args.Error(0)
}

func (sm *StorageMock) Close() {
}
