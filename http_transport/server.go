package http_transport

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"net/http"
	simple_gopher "simple_gopher"
	"time"
)

type Server struct {
	app        *simple_gopher.App
	router     *chi.Mux
	httpServer *http.Server
}

func NewServer(
	app *simple_gopher.App,
) (*Server, error) {
	port := fmt.Sprintf(":%d", app.Config.PORT)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(SecurityMiddleware)
	r.Use(middleware.Timeout(time.Second * 30))
	r.Use(middleware.Heartbeat("/"))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.Config.CorsAllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if app.Config.DebugRoutes {
		r.Use(middleware.Logger)
	}

	// Routing
	swaggerRouter, err := openApi3Router(app.Config)
	if err != nil {
		return nil, err
	}
	r.Route("/docs", swaggerRouter)

	httpServer := &http.Server{
		Addr:              port,
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return &Server{
		app:        app,
		router:     r,
		httpServer: httpServer,
	}, nil
}

func (s *Server) StartAndListen() error {
	log.Info().Msgf("Server started on port :%d", s.app.Config.PORT)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
