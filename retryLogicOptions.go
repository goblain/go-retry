package retry

import (
	"fmt"
	"time"
)

func WithMaxAttempts(count int32) RetryOption {
	return func(rl *RetryLogic) error {
		rl.maxAttempts = count
		return nil
	}
}

func WithLinearBackoff(delay time.Duration) RetryOption {
	return func(rl *RetryLogic) error {
		if delay <= 0 {
			return fmt.Errorf("negative duration not supported")
		}
		rl.baseDelay = delay
		rl.delay = rl.baseDelay
		return nil
	}
}

func WithExponentialBackoff(baseDelay, maxDelay time.Duration, exponentialFactor float32) RetryOption {
	return func(rl *RetryLogic) error {
		if baseDelay <= 0 {
			return fmt.Errorf("negative duration not supported")
		}
		if exponentialFactor <= 1 {
			return fmt.Errorf("exponentialFactor needs to be > 1")
		}
		rl.baseDelay = baseDelay
		rl.delay = rl.baseDelay
		rl.exponentialFactor = exponentialFactor
		rl.maxDelay = maxDelay
		if rl.maxDelay <= 0 {
			rl.maxDelay = 3600 // 1h
		}
		return nil
	}
}
