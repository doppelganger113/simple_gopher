package middleware

import "net/http"

func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-DNS-Prefetch-Control", "off")

		next.ServeHTTP(w, r)
	})
}
