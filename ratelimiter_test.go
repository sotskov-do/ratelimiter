package ratelimiter

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	type args struct {
		limit         int64
		refreshPeriod time.Duration
	}

	tests := []struct {
		name            string
		args            args
		iterations      int
		timeForRequests time.Time
		count           int64
	}{
		{
			name: "inLimit",
			args: args{
				limit:         10,
				refreshPeriod: time.Second,
			},
			iterations:      5,
			timeForRequests: time.Now().Add(5 * time.Second),
			count:           5,
		},
		{
			name: "outOfLimit",
			args: args{
				limit:         10,
				refreshPeriod: time.Second,
			},
			iterations:      6,
			timeForRequests: time.Now().Add(50 * time.Millisecond),
			count:           1,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewRateLimiter(tt.args.limit, tt.args.refreshPeriod)
			var countAcquired int64

			for i := 0; i < tt.iterations; i++ {
				if time.Now().After(tt.timeForRequests) {
					break
				}
				r.Acquire()
				countAcquired++
			}

			if countAcquired != tt.count {
				t.Errorf(
					"Wrong current acquired amount: got %v, want %v Acquire(%v;%v)",
					countAcquired, tt.count, tt.args.limit, tt.args.refreshPeriod,
				)
			}
		})
	}
}
