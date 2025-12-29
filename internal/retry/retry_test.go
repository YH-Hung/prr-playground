package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoSuccess(t *testing.T) {
	callCount := 0
	fn := func() error {
		callCount++
		return nil
	}

	err := Do(context.Background(), 3, fn, func(error) bool { return true })
	if err != nil {
		t.Errorf("Do() should succeed, got error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("Function should be called once, got %d calls", callCount)
	}
}

func TestDoRetryableFailure(t *testing.T) {
	callCount := 0
	retryableErr := errors.New("retryable error")

	fn := func() error {
		callCount++
		if callCount < 3 {
			return retryableErr
		}
		return nil
	}

	isRetryable := func(err error) bool {
		return err == retryableErr
	}

	err := Do(context.Background(), 5, fn, isRetryable)
	if err != nil {
		t.Errorf("Do() should eventually succeed, got error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("Function should be called 3 times, got %d calls", callCount)
	}
}

func TestDoNonRetryableFailure(t *testing.T) {
	callCount := 0
	nonRetryableErr := errors.New("non-retryable error")

	fn := func() error {
		callCount++
		return nonRetryableErr
	}

	isRetryable := func(err error) bool {
		return false
	}

	err := Do(context.Background(), 3, fn, isRetryable)
	if err != nonRetryableErr {
		t.Errorf("Do() should return non-retryable error, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("Function should be called once for non-retryable error, got %d calls", callCount)
	}
}

func TestDoExhaustRetries(t *testing.T) {
	callCount := 0
	retryableErr := errors.New("retryable error")

	fn := func() error {
		callCount++
		return retryableErr
	}

	isRetryable := func(err error) bool {
		return true
	}

	maxRetries := 3
	err := Do(context.Background(), maxRetries, fn, isRetryable)
	if err != retryableErr {
		t.Errorf("Do() should return error after exhausting retries, got: %v", err)
	}
	expectedCalls := maxRetries + 1 // Initial attempt + retries
	if callCount != expectedCalls {
		t.Errorf("Function should be called %d times, got %d calls", expectedCalls, callCount)
	}
}

func TestDoContextCancellation(t *testing.T) {
	callCount := 0
	retryableErr := errors.New("retryable error")

	fn := func() error {
		callCount++
		return retryableErr
	}

	isRetryable := func(err error) bool {
		return true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := Do(ctx, 100, fn, isRetryable)
	if err != context.DeadlineExceeded {
		t.Errorf("Do() should return context error, got: %v", err)
	}
	// Should have made at least one call, but not all 100 retries
	if callCount < 1 || callCount > 10 {
		t.Errorf("Expected 1-10 calls before timeout, got %d calls", callCount)
	}
}
