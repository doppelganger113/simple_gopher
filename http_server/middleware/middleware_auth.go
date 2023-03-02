package middleware

import (
	"api/auth"
	"api/http_server/authenticator"
	"api/http_server/http_util"
	"api/http_server/middleware/keys"
	"context"
	"github.com/rs/zerolog"
	"net/http"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func Authorize(
	logger *zerolog.Logger, validator authenticator.Authenticator, group auth.Role,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			token := http_util.GetTokenFromHeader(authHeader)
			ctx := r.Context()

			isValid, username, err := validator.IsTokenValid(ctx, token, group)
			if err != nil || !isValid {
				if err != nil {
					logger.Warn().Msgf("failed token validation: %s", err)
				}
				http_util.WriteJson(w, http.StatusUnauthorized, http_util.NewFailureResponse("Unauthorized"))
				return
			}

			updatedReq := r.WithContext(
				context.WithValue(ctx, keys.UserAuthDtoKey, auth.AuthorizationDto{
					Header:   authHeader,
					Username: username,
					Role:     group,
				}),
			)
			next.ServeHTTP(w, updatedReq)
		})
	}
}
