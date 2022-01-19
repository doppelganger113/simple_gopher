package concurrency

import (
	"context"
	"time"
)

// SleepSecondsWithContext sleeps for a set period of seconds, but checks if context has ended every second
func SleepSecondsWithContext(ctx context.Context, seconds uint) error {
	var i uint = 0
	for ; i < seconds; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(time.Second)
		}
	}

	return nil
}

type Value []byte

func SlowOperation(ctx context.Context, seconds uint) (Value, error) {
	select {
	case <-time.After(time.Second * time.Duration(seconds)):
		return []byte("Completed hard work."), nil
	case <-ctx.Done():
		return []byte(""), ctx.Err()
	}
}

func Stream(ctx context.Context, out chan<- Value) error {
	derivedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := SlowOperation(derivedCtx, 5)
	if err != nil {
		return err
	}

	for {
		select {
		case out <- res:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
