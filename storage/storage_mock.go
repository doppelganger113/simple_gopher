package storage

import (
	"context"
)

type Mock struct {
}

func (sm Mock) Connect(_ context.Context, _ string) error {
	return nil
}

func (sm Mock) Close() {
}
