package retry

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type AttemptableFuncIterator struct {
	delay            time.Duration
	successIteration int
	currentIteration int
}

func (afi *AttemptableFuncIterator) AttemptableFunc() error {
	afi.currentIteration++
	time.Sleep(afi.delay)
	if afi.currentIteration == afi.successIteration {
		return nil
	}
	return fmt.Errorf("failed iteration")
}

func TestRetryLogic(t *testing.T) {
	type args struct {
		opts []RetryOption
	}
	tests := []struct {
		name         string
		args         args
		iter         *AttemptableFuncIterator
		want         *RetryLogic
		wantErr      bool
		wantDuration time.Duration
		wantAttempts int32
	}{
		{
			name: "linear 15 fail @ 10",
			args: args{
				opts: []RetryOption{
					WithLinearBackoff(time.Second),
				},
			},
			iter: &AttemptableFuncIterator{delay: time.Second, successIteration: 15},
			want: &RetryLogic{
				attempt:           0,
				maxAttempts:       10,
				delay:             time.Second,
				maxDelay:          0,
				baseDelay:         time.Second,
				exponentialFactor: 1,
				autoReset:         false,
			},
			wantErr:      true,
			wantDuration: time.Second * 10,
			wantAttempts: 11,
		},
		{
			name: "linear 5 success",
			args: args{
				opts: []RetryOption{
					WithLinearBackoff(time.Second),
				},
			},
			iter: &AttemptableFuncIterator{delay: time.Second, successIteration: 5},
			want: &RetryLogic{
				attempt:           0,
				maxAttempts:       10,
				delay:             time.Second,
				maxDelay:          0,
				baseDelay:         time.Second,
				exponentialFactor: 1,
				autoReset:         false,
			},
			wantErr:      false,
			wantDuration: time.Second * 5,
			wantAttempts: 5,
		},
		{
			name: "exponential 5 success",
			args: args{
				opts: []RetryOption{
					WithExponentialBackoff(time.Second, time.Second*5, 2),
				},
			},
			iter: &AttemptableFuncIterator{delay: time.Second, successIteration: 5},
			want: &RetryLogic{
				attempt:           0,
				maxAttempts:       10,
				delay:             time.Second,
				maxDelay:          time.Second * 5,
				baseDelay:         time.Second,
				exponentialFactor: 2,
				autoReset:         false,
			},
			wantErr:      false,
			wantDuration: time.Second * 15,
			wantAttempts: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ttp := tt
			t.Parallel()
			rl, _ := NewRetryLogic(tt.args.opts...)
			if !reflect.DeepEqual(rl, tt.want) {
				t.Errorf("NewRetryLogic() got = %v, want %v", rl, tt.want)
			}
			started := time.Now()
			err := rl.ExecuteFunc(ttp.iter.AttemptableFunc)
			duration := time.Now().Sub(started)
			assert.Greater(t, duration, ttp.wantDuration)
			assert.Equal(t, ttp.wantAttempts, rl.attempt)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
		})
	}
}
