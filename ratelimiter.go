package ratelimiter

import (
	"sync/atomic"
	"time"
)

type Limiter interface {
	Acquire()
	SetNewLimit(limit int64, t time.Duration)
}

var _ Limiter = (*RateLimiter)(nil)

type RateLimiter struct {
	state      int64
	padding    [56]byte
	maxSlack   time.Duration
	perRequest time.Duration
}

func NewRateLimiter(limit int64, t time.Duration) *RateLimiter {
	perRequest := t / time.Duration(limit)
	slack := 10 // TODO: to config

	r := &RateLimiter{
		perRequest: perRequest,
		maxSlack:   time.Duration(slack) * perRequest,
	}

	atomic.StoreInt64(&r.state, time.Now().UnixNano())

	return r
}

func (r *RateLimiter) Acquire() {
	var now, next int64

	for {
		now = time.Now().UnixNano()
		last := atomic.LoadInt64(&r.state)

		switch {
		case r.maxSlack > 0 && now-last > int64(r.maxSlack)+int64(r.perRequest):
			// if a lot of time passed between Acquire() calls
			next = now - int64(r.maxSlack)
		default:
			next = last + int64(r.perRequest)
		}

		if atomic.CompareAndSwapInt64(&r.state, last, next) {
			break
		}
	}

	waitTime := time.Duration(next - now)
	if waitTime > 0 {
		time.Sleep(waitTime)
	}
}

func (r *RateLimiter) SetNewLimit(limit int64, t time.Duration) {
	r.perRequest = t / time.Duration(limit)
}
