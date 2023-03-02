package auth

import (
	"api/storage"
	"context"
)

type Mock struct {
}

func (auth *Mock) FetchAndSetKeySet(_ context.Context) error {
	return nil
}
func (auth *Mock) IsTokenValid(
	_ context.Context, _ string, _ Role,
) (valid bool, username string, err error) {
	return false, "", err
}

func (auth *Mock) GetUserAttributes(
	_ context.Context, _ string,
) (UserAttributes, error) {
	return UserAttributes{}, nil
}

func (auth *Mock) GetOrSyncUser(
	_ context.Context, _ AuthorizationDto,
) (storage.User, error) {
	return storage.User{}, nil
}

func (auth *Mock) StartConsumingPostAuthAsync(_ context.Context) {
}
func (auth *Mock) Shutdown() error {
	return nil
}
