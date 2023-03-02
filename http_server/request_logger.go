package http_server

import (
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
	"time"
)

func CreateRequestLogger(logger *zerolog.Logger, config Config) func(http.Handler) http.Handler {
	log := logger.With().
		Timestamp().
		Str("service", "api").
		Str("host", config.Domain).
		Logger()
	c := alice.New()

	// Install the logger handler with default output on the console
	c = c.Append(hlog.NewHandler(log))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to that handler, all our logs will come with some pre-populated fields.
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("Request")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	return func(next http.Handler) http.Handler {
		// Here is your final handler
		return c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the logger from the request's context. You can safely assume it
			// will be always there: if the handler is removed, hlog.FromRequest
			// will return a no-op logger.
			//hlog.FromRequest(r).Info().
			//	Str("method", r.Method).
			//	Str("uri", r.RequestURI).
			//	Msg("Request")

			next.ServeHTTP(w, r)
		}))
	}
}
