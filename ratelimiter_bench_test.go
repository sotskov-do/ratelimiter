package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkRateLimiter_Acquire(b *testing.B) {
	r := NewRateLimiter(100, 5*time.Second)

	wg := &sync.WaitGroup{}
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			r.Acquire()
		}()
	}

	wg.Wait()
}
