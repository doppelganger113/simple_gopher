package http_transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"simple_gopher/auth"
)

type ContextKey string

var (
	UserAuthDtoKey ContextKey = "UserAuthDtoKey"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func Authorize(
	next http.HandlerFunc, validator Authenticator, group auth.Role,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := GetTokenFromHeader(authHeader)
		ctx := r.Context()

		isValid, username, err := validator.IsTokenValid(ctx, token, group)
		if err != nil || !isValid {
			if err != nil {
				log.Println(fmt.Errorf("failed token validation %w", err))
			}
			respondJson(w, http.StatusUnauthorized, newFailureResponse("Unauthorized"))
			return
		}

		updatedReq := r.WithContext(
			context.WithValue(ctx, UserAuthDtoKey, auth.AuthorizationDto{
				Header:   authHeader,
				Username: username,
				Role:     group,
			}),
		)
		next(w, updatedReq)
	}
}
