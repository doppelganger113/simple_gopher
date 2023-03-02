package concurrency

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNewRetry(t *testing.T) {
	newRetry := NewRetry(2, 3*time.Second)

	if !reflect.DeepEqual(newRetry, &Retry{
		retries:        2,
		delay:          3 * time.Second,
		baseBackoff:    time.Second,
		maximumBackoff: time.Minute,
	}) {
		t.Fatal("expected retry to be properly set")
	}
}

func TestRetry_Execute(t *testing.T) {
	newRetry := Retry{
		retries:        3,
		delay:          3 * time.Microsecond,
		baseBackoff:    time.Microsecond,
		maximumBackoff: 10 * time.Microsecond,
	}

	counter := 0
	var executionCounts []uint

	err := newRetry.Execute(context.Background(), func(_ context.Context, retryCount uint) error {
		executionCounts = append(executionCounts, retryCount)

		if counter++; counter < 3 {
			return errors.New("an error")
		}

		return nil
	})
	if err != nil {
		t.Fatal("failed execute", err)
	}

	if counter != 3 {
		t.Fatalf("Expected counter to be 3, got %d", counter)
	}

	if !reflect.DeepEqual(executionCounts, []uint{0, 1, 2}) {
		t.Fatal("expected retry counts to be equal")
	}
}

func TestRetry_Execute_NoRetry(t *testing.T) {
	newRetry := Retry{
		retries:        3,
		delay:          3 * time.Microsecond,
		baseBackoff:    time.Microsecond,
		maximumBackoff: 10 * time.Microsecond,
	}

	counter := 0
	err := newRetry.Execute(context.Background(), func(_ context.Context, retryCount uint) error {
		counter++
		return nil
	})
	if err != nil {
		t.Fatal("expected no errors")
	}
	if counter != 1 {
		t.Fatal("expected only a single execution")
	}
}

func TestRetry_Execute_Err(t *testing.T) {
	newRetry := Retry{
		retries:        3,
		delay:          3 * time.Microsecond,
		baseBackoff:    time.Microsecond,
		maximumBackoff: 10 * time.Microsecond,
	}

	counter := 0
	var executionCounts []uint

	err := newRetry.Execute(context.Background(), func(_ context.Context, retryCount uint) error {
		executionCounts = append(executionCounts, retryCount)

		if counter++; counter < 5 {
			return errors.New("an error")
		}

		return nil
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !reflect.DeepEqual(executionCounts, []uint{0, 1, 2, 3}) {
		t.Fatal("didn't receive expected execution counts")
	}
}
