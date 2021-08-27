package cloud_patterns

import (
	"context"
	"math/rand"
	"time"
)

type RetryEffector func(ctx context.Context, retryCount uint) error

type retry struct {
	retries uint
	delay   time.Duration
}

func NewRetry(delay time.Duration) *retry {
	return &retry{delay: delay}
}

// Execute will retry failed request with a jitter backoff algorithm. Please
// execute a rand.Seed(time.Now().UTC().UnixNano()) at the start of the main
// function to have different numbers generated
func (r *retry) Execute(ctx context.Context, effector RetryEffector) error {
	err := effector(ctx, 0)
	baseBackoff, maximumBackoff := time.Second, time.Minute

	var retryCount uint
	for backoff := baseBackoff; err != nil; backoff <<= 1 {
		retryCount++

		if backoff > maximumBackoff {
			backoff = maximumBackoff
		}

		jitter := rand.Int63n(int64(backoff * 3))
		sleep := baseBackoff + time.Duration(jitter)
		time.Sleep(sleep)
		err = effector(ctx, retryCount)
	}

	return err
}
