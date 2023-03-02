package api

import (
	"api/core"
	"api/http_server"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

var (
	GitCommit string
	BuildDate string
	Version   string
)

// Bootstrap - starts the http server and returns channel that signals termination
func Bootstrap(app *core.App, logger *zerolog.Logger, configurations ...Configure) (<-chan error, error) {
	closed := make(chan error)
	bootstrapConfig := newConfiguration(configurations...)

	logger.Info().Msgf("[Build]: Version %s BuildDate %s GitCommit %s", Version, BuildDate, GitCommit)

	appCtx := context.Background()
	defer appCtx.Done()

	errChannel := make(chan error)

	appInitTimeout, appInitCancel := context.WithTimeout(
		appCtx, bootstrapConfig.initTimeout,
	)
	defer appInitCancel()

	logger.Info().Msg("Initializing app...")

	if err := app.Init(appInitTimeout, appCtx); err != nil {
		return nil, fmt.Errorf("failed initializing app: %w", err)
	}

	server, err := http_server.StartNewConfiguredAndListenChannel(logger, app, errChannel)
	if err != nil {
		return nil, fmt.Errorf("failed starting the server: %w", err)
	}
	go bootstrapConfig.interrupt(errChannel)

	go func() {
		fatalErr := <-errChannel

		shutdownGracefully(app, logger, server, bootstrapConfig.shutdownTimeout)
		closed <- fatalErr
	}()

	return closed, nil
}

func shutdownGracefully(app *core.App, logger *zerolog.Logger, server *http_server.Server, timeout time.Duration) {
	if app == nil && server == nil {
		return
	}
	logger.Info().Msg("Gracefully shutting down...")

	gracefullCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if server != nil {
		if err := server.Shutdown(gracefullCtx); err != nil {
			logger.Error().Msgf("Error shutting down the server: %s", err.Error())
		} else {
			logger.Info().Msg("HttpServer gracefully shut down")
		}
	}

	if app != nil {
		if err := app.Shutdown(gracefullCtx); err != nil {
			logger.Error().Msg(err.Error())
		} else {
			logger.Info().Msg("application shut down successfully")
		}
	}
}
