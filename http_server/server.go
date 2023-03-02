package http_server

import (
	"api/core"
	coremiddleware "api/http_server/middleware"
	"api/http_server/openapi"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Server struct {
	Port       uint
	router     *chi.Mux
	httpServer *http.Server
	logger     *zerolog.Logger
}

func NewServer(logger *zerolog.Logger, config Config, app *core.App) (*Server, error) {
	port := fmt.Sprintf(":%d", config.Port)

	r := chi.NewRouter()
	r.Use(CreateRequestLogger(logger, config))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(coremiddleware.Secure)
	r.Use(middleware.Timeout(config.Timeout))
	r.Use(middleware.Heartbeat(config.HeartbeatUrl))
	r.Use(Cors(config.CorsAllowOrigins))

	// Enable httprate request limiter of 100 requests per minute.
	//
	// In the code example below, rate-limiting is bound to the request IP address
	// via the LimitByIP middleware handler.
	//
	// To have a single rate-limiter for all requests, use httprate.LimitAll(..).
	//
	// Please see _example/main.go for other more, or read the library code.
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// Routing
	swaggerRouter, err := openapi.NewOpenApi3Router(openapi.Config{
		BasicAuthUsername:          config.BasicAuthUsername,
		BasicAuthPassword:          config.BasicAuthPassword,
		BasicAuthRealm:             config.BasicAuthRealm,
		Domain:                     config.Domain,
		OAuth2TokenUrl:             config.OAuth2TokenUrl,
		OAuth2AuthorizationCodeUrl: config.OAuth2AuthorizationCodeUrl,
	})
	if err != nil {
		return nil, err
	}
	r.Route("/docs", swaggerRouter)

	imagesHandler := NewImageHandler(logger, app.Auth, app.ImagesService)
	r.Route("/api/v1/images", imagesHandler.CreateRouter())

	httpServer := &http.Server{
		Addr:              port,
		Handler:           r,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
	}

	return &Server{
		router:     r,
		httpServer: httpServer,
		Port:       config.Port,
		logger:     logger,
	}, nil
}

// StartNewConfiguredAndListenChannel boots configuration, creates and starts the server with
// err channel which is used to signal when the server closes
func StartNewConfiguredAndListenChannel(
	logger *zerolog.Logger, app *core.App, errChannel chan<- error,
) (*Server, error) {
	var server *Server

	httpConfig := NewDefaultConfig()
	if err := httpConfig.LoadFromEnv(); err != nil {
		return nil, err
	}

	server, err := NewServer(logger, httpConfig, app)
	if err != nil {
		return nil, err
	}

	go func() {
		errChannel <- server.StartAndListen()
	}()

	return server, nil
}

func (s *Server) StartAndListen() error {
	s.logger.Info().Msgf("Server started on port :%d", s.Port)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
