# RateLimiter

This package provides a Golang implementation of the token-bucket rate limit algorithm. This implementation refills the bucket based on specified time interval.

Create a rate limiter with a maximum number of operations to perform per specified time interval. Call `Acquire()` before each operation. `Acquire` will return error if out of limit. In error message you can check time to bucket refilling. Parameters of rate limiter can be change after creating by calling `SetLimit(limit int64)` and `SetRefreshPeriod(t time.Duration)`


### Example 
``` golang
totalOperations := 10
timeBetweenOperations := 1 * time.Second
limit := 5 // amount of operations can be completed during refresh period
refreshPeriod := 10 * time.Second // interval of time during which a limited number of operations are performed
limiter := ratelimiter.NewRateLimiter(limit, refreshPeriod)

for i := 0; i < totalOperations; i++ {
    err := limiter.Acquire()
    if err != nil {
        fmt.Printf("error during operation #%v: %v\n", i+1, err.Error())
    }
    time.Sleep(timeBetweenOperations)
}

// Output:
// error during operation #6: can't do the operation, try again after 4.999226689s
// error during operation #7: can't do the operation, try again after 3.998541803s
// error during operation #8: can't do the operation, try again after 2.998311454s
// error during operation #9: can't do the operation, try again after 1.997595771s
// error during operation #10: can't do the operation, try again after 997.332071ms
```