package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"simple_gopher/http_transport"
	"simple_gopher/simple_gopher"
	"syscall"
	"time"
)

var (
	GitCommit string
	BuildDate string
	Version   string
)

// Shutdown timeout should be as long or more as the request timeout
const (
	shutdownTimeoutSeconds = 30
	appInitTimeoutSeconds  = 30
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msgf("[Build info]: Version %s BuildDate %s GitCommit %s", Version, BuildDate, GitCommit)

	appCtx := context.Background()
	defer appCtx.Done()

	errChannel := make(chan error)

	log.Info().Msg("Loading configuration...")

	// Configurations
	config, err := simple_gopher.NewConfigFromEnv()
	if err != nil {
		log.Fatal().Msgf("failed creating configuration from env: %s", err.Error())
		return
	}

	app, buildErr := simple_gopher.CreateApp(config)
	if buildErr != nil {
		log.Fatal().Msgf("failed creating app: %s", buildErr.Error())
		return
	}

	appInitTimeout, appInitCancel := context.WithTimeout(
		appCtx, appInitTimeoutSeconds*time.Second,
	)
	defer appInitCancel()

	log.Info().Msg("Initializing app...")

	if err = app.Init(appInitTimeout, appCtx); err != nil {
		log.Fatal().Msgf("failed initializing app: %s", err.Error())
		return
	}

	server, err := http_transport.StartNewConfiguredAndListenChannel(http_transport.Handlers{
		ImagesHandler: app.ImagesService,
		Authenticator: app.Auth,
	}, errChannel)
	if err != nil {
		log.Fatal().Msgf("failed starting the server: %s", err.Error())
		return
	}
	go listenForInterrupt(errChannel)

	fatalErr := <-errChannel
	log.Info().Msgf("Closing server: %s", fatalErr.Error())

	shutdownGracefully(app, server, shutdownTimeoutSeconds*time.Second)
}

func listenForInterrupt(errChannel chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	errChannel <- fmt.Errorf("%s", <-c)
}

func shutdownGracefully(app *simple_gopher.App, server *http_transport.Server, timeout time.Duration) {
	if app == nil && server == nil {
		return
	}
	log.Info().Msg("Gracefully shutting down...")

	gracefullCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if server != nil {
		if err := server.Shutdown(gracefullCtx); err != nil {
			log.Error().Msgf("Error shutting down the server: %s", err.Error())
		} else {
			log.Info().Msg("HttpServer gracefully shut down")
		}
	}

	if app != nil {
		if err := app.Shutdown(gracefullCtx); err != nil {
			log.Error().Msg(err.Error())
		} else {
			log.Info().Msg("application shut down successfully")
		}
	}
}
