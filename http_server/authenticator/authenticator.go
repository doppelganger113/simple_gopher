package authenticator

import (
	"api/auth"
	"context"
)

type Authenticator interface {
	IsTokenValid(
		ctx context.Context, tokenString string, role auth.Role,
	) (isValid bool, token string, err error)
}
