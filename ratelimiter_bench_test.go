package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkRateLimiter_Acquire(b *testing.B) {
	r := NewRateLimiter(10000, 5*time.Second)
	wg := &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Acquire()
		}()
	}

	wg.Wait()
}
