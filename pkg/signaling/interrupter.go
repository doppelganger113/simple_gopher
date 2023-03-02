package signaling

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type EventEmitter func(err chan<- error)

// TerminateEventEmitter - sends an error on termination signal, eg ctrl+c
func TerminateEventEmitter(err chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	err <- fmt.Errorf("%s", <-c)
}

// Forward - forwards error from one channel to another
func Forward(from <-chan error) func(chan<- error) {
	return func(to chan<- error) {
		to <- fmt.Errorf("%s", <-from)
	}
}
