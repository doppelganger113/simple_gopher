package cloud_patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ServiceUnreachableErr = errors.New("service unreachable")

type Circuit func(context.Context) (Value, error)

type CircuitBreaker struct {
	sync.RWMutex
	failureThreshold    uint
	consecutiveFailures uint
	lastAttempt         time.Time
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{}
}

func (cb *CircuitBreaker) Run(ctx context.Context, circuit Circuit) (Value, error) {
	if retry := cb.shouldRetry(); !retry {
		return nil, ServiceUnreachableErr
	}

	response, err := circuit(ctx)

	return cb.handleResponse(response, err)
}

func (cb *CircuitBreaker) handleResponse(response Value, err error) (Value, error) {
	cb.Lock()
	defer cb.Unlock()

	cb.lastAttempt = time.Now()
	if err != nil {
		cb.consecutiveFailures++
		return response, err
	}

	cb.consecutiveFailures = 0

	return response, nil
}

func (cb *CircuitBreaker) shouldRetry() bool {
	cb.RLock()
	defer cb.RUnlock()
	diff := cb.consecutiveFailures - cb.failureThreshold

	if diff >= 0 {
		shouldRetryAt := cb.lastAttempt.Add(time.Second * 2 << diff)
		if !time.Now().After(shouldRetryAt) {
			return false
		}
	}

	return true
}
