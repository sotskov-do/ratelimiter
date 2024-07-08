package ratelimiter

import (
	"sync/atomic"
	"time"
)

type Limiter interface {
	Acquire() error
	SetLimit(limit int64)
	SetRefreshPeriod(t time.Duration)
	GetCurrentAmount() int64
	GetLastUpdate() time.Time
}

type RateLimiter struct {
	currentAmount int64
	lastUpdate    int64
	limit         int64
	refreshPeriod time.Duration
}

func NewRateLimiter(limit int64, t time.Duration) *RateLimiter {
	r := &RateLimiter{
		limit:         limit,
		refreshPeriod: t,
	}

	atomic.StoreInt64(&r.currentAmount, limit)
	atomic.StoreInt64(&r.lastUpdate, time.Now().UnixNano())

	return r
}

func (r *RateLimiter) Acquire() error {
	now := time.Now().UnixNano()
	lastUpdate := atomic.LoadInt64(&r.lastUpdate)
	timeElapsed := now - lastUpdate

	currentAmount := atomic.LoadInt64(&r.currentAmount)

	if timeElapsed > r.refreshPeriod.Nanoseconds() {
		currentAmount = r.limit
		atomic.StoreInt64(&r.lastUpdate, now)
	}

	currentAmount--
	if currentAmount < 0 {
		waitTime := r.refreshPeriod.Nanoseconds() - timeElapsed
		time.Sleep(time.Duration(waitTime))
	}

	atomic.StoreInt64(&r.currentAmount, currentAmount)

	return nil
}

func (r *RateLimiter) SetLimit(limit int64) {
	r.limit = limit
}

func (r *RateLimiter) SetRefreshPeriod(t time.Duration) {
	r.refreshPeriod = t
}

func (r *RateLimiter) GetCurrentAmount() int64 {
	return atomic.LoadInt64(&r.currentAmount)
}

func (r *RateLimiter) GetLastUpdate() time.Time {
	return time.Unix(0, atomic.LoadInt64(&r.lastUpdate))
}
