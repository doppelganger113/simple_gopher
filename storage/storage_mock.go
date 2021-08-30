package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (sm *Mock) Connect(ctx context.Context, connectionUrl string) error {
	args := sm.Called(ctx, connectionUrl)
	return args.Error(0)
}

func (sm *Mock) Close() {
}
