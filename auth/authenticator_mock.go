package auth

import (
	"context"
	"github.com/stretchr/testify/mock"
	"simple_gopher/storage"
)

type AuthMock struct {
	mock.Mock
}

func (auth *AuthMock) FetchAndSetKeySet(ctx context.Context) error {
	args := auth.Called(ctx)

	return args.Error(0)
}
func (auth *AuthMock) IsTokenValid(
	ctx context.Context, token string, authGroup Role,
) (valid bool, username string, err error) {
	args := auth.Called(ctx, token, authGroup)

	return args.Bool(0), args.String(1), args.Error(2)
}

func (auth *AuthMock) GetUserAttributes(
	ctx context.Context, username string,
) (*UserAttributes, error) {
	args := auth.Called(ctx, username)

	return args.Get(0).(*UserAttributes), args.Error(1)
}

func (auth *AuthMock) GetOrSyncUser(
	ctx context.Context, authorization AuthorizationDto,
) (*storage.User, error) {
	args := auth.Called(ctx, authorization)

	return args.Get(0).(*storage.User), args.Error(1)
}
func (auth *AuthMock) StartConsumingPostAuthAsync(ctx context.Context) {
	_ = auth.Called(ctx)
}
func (auth *AuthMock) Shutdown() error {
	args := auth.Called()

	return args.Error(0)
}
