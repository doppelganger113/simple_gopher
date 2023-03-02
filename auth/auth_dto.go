package auth

import (
	"context"
)

type AuthorizationDto struct {
	Header   string
	Username string
	Role     Role
}

func ExtractAuthorizationDto(ctx context.Context, key interface{}) (AuthorizationDto, error) {
	value := ctx.Value(key)
	if value == nil {
		return AuthorizationDto{}, ErrMissingAuthDto
	}
	authDto, ok := value.(AuthorizationDto)
	if !ok {
		return AuthorizationDto{}, ErrMissingAuthDto
	}

	return authDto, nil
}
