package retry

import (
	"time"
)

const (
	DEFAULT_MAX_ATTEMPTS = 10
	DEFAULT_DELAY        = time.Second
)

func NewRetryLogic(opts ...RetryOption) (*RetryLogic, error) {
	rl := &RetryLogic{
		attempt:           0,
		maxAttempts:       DEFAULT_MAX_ATTEMPTS,
		baseDelay:         DEFAULT_DELAY,
		delay:             DEFAULT_DELAY,
		exponentialFactor: 1,
	}
	for _, opt := range opts {
		if err := opt(rl); err != nil {
			return nil, err
		}
	}
	return rl, nil
}

// Copy returns a goroutine safe instance of the logic
func (rl *RetryLogic) Copy() *RetryLogic {
	newRL := &RetryLogic{}
	*newRL = *rl
	return newRL
}

func (rl *RetryLogic) Reset() {
	rl.attempt = 0
	rl.delay = rl.baseDelay
}

// Attempt returns false if you're not allowed to perform another execution
// otherwise, if there are still executions eligible, it blocks until
// an execution can happen according to backoff and return true
func (rl *RetryLogic) Attempt() bool {
	rl.attempt++
	if rl.attempt == 1 {
		rl.lastAttemptTime = time.Now()
		return true
	}
	if rl.attempt > rl.maxAttempts {
		if rl.autoReset {
			rl.attempt = 0
			rl.delay = rl.baseDelay
		}
		return false
	}

	time.Sleep(rl.delay)
	rl.delay = time.Duration(float32(rl.delay) * rl.exponentialFactor)
	if rl.delay > rl.maxDelay {
		rl.delay = rl.maxDelay
	}

	rl.lastAttemptTime = time.Now()
	return true
}

func (rl *RetryLogic) AttemptDone() {
	rl.lastAttemptTime = time.Now()
}

// ExecuteFunc retries any AttemptableFunc untill it is successful
// or no longer allowed to retry and returns the last error
func (r *RetryLogic) ExecuteFunc(f AttemptableFunc) error {
	var err error
	for r.Attempt() {
		err = f()
		if err != nil {
			continue
		}
		break
	}
	return err
}

// ExecuteFuncI retries any AttemptableFuncI until it is successful
// or no longer allowed to retry and returns a single output interface and last error
func (r *RetryLogic) ExecuteFuncI(f AttemptableFuncI) (interface{}, error) {
	var err error
	var output interface{}
	for r.Attempt() {
		output, err = f()
		if err != nil {
			continue
		}
		break
	}
	return output, err
}
