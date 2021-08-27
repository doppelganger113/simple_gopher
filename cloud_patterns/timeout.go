package cloud_patterns

import "context"

type SlowFunction func(ctx context.Context) (Value, error)

func Timeout(ctx context.Context, fn SlowFunction) (Value, error) {
	chres := make(chan Value)
	cherr := make(chan error)

	go func() {
		res, err := fn(ctx)
		chres <- res
		cherr <- err
	}()

	select {
	case res := <-chres:
		return res, <-cherr
	case <-ctx.Done():
		return Value(""), ctx.Err()
	}
}
