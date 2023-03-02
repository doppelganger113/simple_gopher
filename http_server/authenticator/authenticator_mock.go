package authenticator

import (
	"api/auth"
	"context"
)

const bearerTokenMock = "tokenMock"

type Mock struct{}

func (a Mock) IsTokenValid(
	_ context.Context, tokenString string, _ auth.Role,
) (isValid bool, token string, err error) {
	if tokenString == bearerTokenMock {
		return true, bearerTokenMock, nil
	}

	return false, "", nil
}
