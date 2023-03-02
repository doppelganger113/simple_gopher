package main

import (
	"api"
	"api/logger"
)

func main() {
	log := logger.NewLogger()
	app, err := api.InitializeApp(log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed initializing app")
	}
	errCh, err := api.Bootstrap(app, log)
	if err != nil {
		log.Error().Err(err).Msgf("failed starting application")
	}
	log.Err(<-errCh).Msg("Shutdown")
}
