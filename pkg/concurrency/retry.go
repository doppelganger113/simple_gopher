package concurrency

import (
	"context"
	"math/rand"
	"time"
)

type RetryEffector func(ctx context.Context, retryCount uint) error

type Retry struct {
	retries        uint
	delay          time.Duration
	baseBackoff    time.Duration
	maximumBackoff time.Duration
}

func NewRetry(retries uint, delay time.Duration) *Retry {
	return &Retry{
		retries:        retries,
		delay:          delay,
		baseBackoff:    time.Second,
		maximumBackoff: time.Minute,
	}
}

// Execute will retry failed request with a jitter backoff algorithm. Please
// execute a rand.Seed(time.Now().UTC().UnixNano()) at the start of the main
// function to have different numbers generated
func (r *Retry) Execute(ctx context.Context, effector RetryEffector) error {
	err := effector(ctx, 0)

	var retryCount uint
	for backoff := r.baseBackoff; err != nil && retryCount < r.retries; backoff <<= 1 {
		retryCount++

		if backoff > r.maximumBackoff {
			backoff = r.maximumBackoff
		}

		jitter := rand.Int63n(int64(backoff * 3))
		sleep := r.baseBackoff + time.Duration(jitter)
		time.Sleep(sleep)
		err = effector(ctx, retryCount)
	}

	return err
}
