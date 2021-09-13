package http_transport

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Server struct {
	Port       uint
	router     *chi.Mux
	httpServer *http.Server
}

type Handlers struct {
	ImagesHandler ImagesHandler
	Authenticator Authenticator
}

func NewServer(config Config, handlers Handlers) (*Server, error) {
	port := fmt.Sprintf(":%d", config.Port)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(SecurityMiddleware)
	r.Use(middleware.Timeout(config.Timeout))
	r.Use(middleware.Heartbeat(config.HeartbeatUrl))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.CorsAllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if config.DebugRoutes {
		r.Use(middleware.Logger)
	}

	// Routing
	swaggerRouter, err := openApi3Router(config)
	if err != nil {
		return nil, err
	}
	r.Route("/docs", swaggerRouter)
	r.Route("/api/v1/images", ImagesRouter(
		handlers.ImagesHandler,
		handlers.Authenticator,
	))

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
	}, nil
}

// StartNewConfiguredAndListenChannel boots configuration, creates and starts the server with
// err channel which is used to signal when the server closes
func StartNewConfiguredAndListenChannel(
	handlers Handlers, errChannel chan<- error,
) (*Server, error) {
	var server *Server

	httpConfig := NewDefaultConfig()
	if err := httpConfig.LoadFromEnv(); err != nil {
		return nil, err
	}

	server, err := NewServer(httpConfig, handlers)
	if err != nil {
		return nil, err
	}

	go func() {
		errChannel <- server.StartAndListen()
	}()

	return server, nil
}

func (s *Server) StartAndListen() error {
	log.Info().Msgf("Server started on port :%d", s.Port)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
