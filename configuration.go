package api

import (
	"api/pkg/signaling"
	"time"
)

type configuration struct {
	interrupt       signaling.EventEmitter
	shutdownTimeout time.Duration
	initTimeout     time.Duration
}

type Configure func(config *configuration)

func newConfiguration(configurations ...Configure) configuration {
	defaultConfig := configuration{
		shutdownTimeout: 30 * time.Second,
		initTimeout:     30 * time.Second,
		interrupt:       signaling.TerminateEventEmitter,
	}

	for _, configure := range configurations {
		configure(&defaultConfig)
	}

	return defaultConfig
}

func WithInterrupt(interrupter signaling.EventEmitter) Configure {
	return func(config *configuration) {
		config.interrupt = interrupter
	}
}

func WithShutdownTimeout(timeout time.Duration) Configure {
	return func(config *configuration) {
		config.shutdownTimeout = timeout
	}
}

func WithInitTimeout(timeout time.Duration) Configure {
	return func(config *configuration) {
		config.initTimeout = timeout
	}
}
