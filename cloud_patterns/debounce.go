package cloud_patterns

import (
	"context"
	"sync"
	"time"
)

type debounceFirst struct {
	sync.Mutex
	threshold time.Time
	result    Value
	err       error
}

func NewDebounceFirst() *debounceFirst {
	return &debounceFirst{}
}

func (d *debounceFirst) Run(ctx context.Context, circuit Circuit) (Value, error) {
	d.Lock()
	defer func() {
		d.threshold = time.Now()
		d.Unlock()
	}()

	if time.Now().Before(d.threshold) {
		return d.result, d.err
	}

	return circuit(ctx)
}

type debounceLast struct {
	threshold time.Time
	ticker    *time.Ticker
	duration  time.Duration
	result    []byte
	err       error
	once      sync.Once
	m         sync.Mutex
}

func NewDebounceLast(duration time.Duration) *debounceLast {
	return &debounceLast{threshold: time.Now(), duration: duration}
}

func (d *debounceLast) Run(ctx context.Context, circuit Circuit) (Value, error) {
	d.m.Lock()
	defer d.m.Unlock()

	d.threshold = time.Now().Add(d.duration)

	d.once.Do(func() {
		d.ticker = time.NewTicker(time.Millisecond * 100)

		go func() {
			defer d.stopTicker()

			for {
				select {
				case <-d.ticker.C:
					d.execute(ctx, circuit)
				case <-ctx.Done():
					d.handleContextDone(ctx)
				}
			}
		}()
	})

	return d.result, d.err
}

func (d *debounceLast) stopTicker() {
	d.m.Lock()
	defer d.m.Unlock()
	d.ticker.Stop()
	d.once = sync.Once{}
}

func (d *debounceLast) execute(ctx context.Context, circuit Circuit) {
	d.m.Lock()
	defer d.m.Unlock()
	if time.Now().After(d.threshold) {
		d.result, d.err = circuit(ctx)
	}
}

func (d *debounceLast) handleContextDone(ctx context.Context) {
	d.m.Lock()
	defer d.m.Unlock()
	d.result, d.err = Value(""), ctx.Err()
}
