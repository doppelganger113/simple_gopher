package concurrency

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrTooManyCalls = errors.New("too many calls")

type Effector func(ctx context.Context)

type throttle struct {
	tokens   uint
	max      uint
	refill   uint
	duration time.Duration
	once     sync.Once
}

func NewThrottle() *throttle {
	return &throttle{}
}

// Execute is missing goroutine safety
func (t *throttle) Execute(ctx context.Context, effector Effector) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	t.once.Do(func() {
		ticket := time.NewTicker(t.duration)

		go func() {
			defer ticket.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticket.C:
					tokens := t.tokens + t.refill
					if tokens > t.max {
						tokens = t.max
					}
					t.tokens = tokens
				}
			}
		}()
	})

	if t.tokens <= 0 {
		return ErrTooManyCalls
	}

	t.tokens--

	effector(ctx)

	return nil
}
