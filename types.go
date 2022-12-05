package retry

import "time"

type RetryLogic struct {
	attempt           int32
	lastAttemptTime   time.Time
	maxAttempts       int32
	delay             time.Duration
	maxDelay          time.Duration
	baseDelay         time.Duration
	exponentialFactor float32
	autoReset         bool
}

type RetryOption func(rl *RetryLogic) error

type AttemptableFunc func() error
type AttemptableFuncI func() (interface{}, error)
