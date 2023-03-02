package auth

import (
	"api/storage"
	"context"
)

type UserAttributes struct {
	Username string
	Sub      string
	Name     string
	Email    string
}

type Authenticator interface {
	FetchAndSetKeySet(ctx context.Context) error
	IsTokenValid(
		ctx context.Context, tokenString string, requiredGroup Role,
	) (valid bool, username string, err error)
	GetUserAttributes(ctx context.Context, username string) (UserAttributes, error)
	GetOrSyncUser(
		ctx context.Context, authorization AuthorizationDto,
	) (storage.User, error)
	StartConsumingPostAuthAsync(ctx context.Context)
	Shutdown() error
}
