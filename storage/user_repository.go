package storage

import (
	"context"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, dto UserCreationDto) (User, error)
}
