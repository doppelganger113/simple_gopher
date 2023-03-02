package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

func NewLogger(options ...OptionFn) *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	defaultOptions := Options{}
	for _, o := range options {
		o(&defaultOptions)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if defaultOptions.pretty {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		logger = logger.With().
			Str("version", defaultOptions.version).
			Int("pid", os.Getpid()).
			Logger()
	}

	return &logger
}
