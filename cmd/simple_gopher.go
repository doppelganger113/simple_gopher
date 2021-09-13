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

// Shutdown timeout should be as long or more as the request timeout
const (
	shutdownTimeoutSeconds = 30
	appInitTimeoutSeconds  = 30
)

func main() {
	appCtx := context.Background()
	defer appCtx.Done()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting api...")

	errChannel := make(chan error)

	app, buildErr := simple_gopher.InitializeApp()
	if buildErr != nil {
		log.Fatal().Msg(buildErr.Error())
		return
	}
	appInitTimeout, appInitCancel := context.WithTimeout(
		appCtx, appInitTimeoutSeconds*time.Second,
	)
	defer appInitCancel()
	err := app.Init(appInitTimeout, appCtx)
	if err != nil {
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
	log.Info().Msg(fatalErr.Error())

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
