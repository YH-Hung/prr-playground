// Package retry provides utilities for retrying operations with exponential backoff.
package retry

import (
	"context"
	"time"
)

const (
	// DefaultBaseDelay is the initial delay before the first retry.
	DefaultBaseDelay = 50 * time.Millisecond
	// DefaultMultiplier is the factor by which the delay increases after each retry.
	DefaultMultiplier = 2.0
)

// Do executes the given function with retry logic using exponential backoff.
// It retries up to maxRetries times if the error is retryable according to the isRetryable function.
// The context can be used to cancel the operation.
func Do(ctx context.Context, maxRetries int, fn func() error, isRetryable func(error) bool) error {
	var err error
	delay := DefaultBaseDelay

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Try the operation
		err = fn()
		if err == nil {
			return nil
		}

		// Check if we should retry
		if !isRetryable(err) {
			return err
		}

		// Don't wait after the last attempt
		if attempt == maxRetries {
			return err
		}

		// Wait with exponential backoff
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			delay = time.Duration(float64(delay) * DefaultMultiplier)
		}
	}

	return err
}
