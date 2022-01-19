package http_transport

import (
	"context"
	"simple_gopher/auth"
)

const bearerTokenMock = "tokenMock"

type AuthenticatorMock struct{}

func (a AuthenticatorMock) IsTokenValid(
	_ context.Context, tokenString string, _ auth.Role,
) (isValid bool, token string, err error) {
	if tokenString == bearerTokenMock {
		return true, bearerTokenMock, nil
	}

	return false, "", nil
}
