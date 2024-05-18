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
		name                string
		args                args
		iterations          int
		timeBetweenRequests time.Duration
		count               int64
		wantErr             bool
	}{
		{
			name: "inLimit",
			args: args{
				limit:         10,
				refreshPeriod: time.Hour,
			},
			iterations:          5,
			timeBetweenRequests: 250 * time.Millisecond,
			count:               5,
		},
		{
			name: "outOfLimit",
			args: args{
				limit:         5,
				refreshPeriod: time.Hour,
			},
			iterations:          6,
			timeBetweenRequests: 250 * time.Millisecond,
			wantErr:             true,
			count:               0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var err error

			r := NewRateLimiter(tt.args.limit, tt.args.refreshPeriod)

			for i := 0; i < tt.iterations; i++ {
				err = r.Acquire()
				time.Sleep(tt.timeBetweenRequests)
			}

			if tt.wantErr && err == nil {
				t.Errorf("Expected error: Acquire(%v;%v)", tt.args.limit, tt.args.refreshPeriod)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Not expected error: Acquire(%v;%v)", tt.args.limit, tt.args.refreshPeriod)
			}

			if r.GetCurrentAmount() != tt.count {
				t.Errorf(
					"Wrong current limit amount: got %v, want %v Acquire(%v;%v)",
					r.GetCurrentAmount(), tt.count, tt.args.limit, tt.args.refreshPeriod,
				)
			}
		})
	}
}
