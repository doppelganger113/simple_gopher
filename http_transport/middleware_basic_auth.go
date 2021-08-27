package http_transport

import (
	"crypto/subtle"
	"net/http"
)

// newBasicAuth creates a basic auth middleware function
func newBasicAuth(username, password, realm string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			user, pass, ok := req.BasicAuth()

			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				_, _ = w.Write([]byte("Unauthorised.\n"))
				return
			}

			next.ServeHTTP(w, req)
		}
	}
}
