package storage

import "context"

type Storage interface {
	Connect(context.Context, string) error
	Close()
}
