package rate_limiter

import "time"

type slidingWindowRatelimiter struct {
	maxRequests int
	windowSize  time.Duration
	requests    map[string][]time.Time
}

func NewSlidingWindowRatelimiter(maxRequests int, windowSize time.Duration) *slidingWindowRatelimiter {
	return &slidingWindowRatelimiter{
		maxRequests: maxRequests,
		windowSize:  windowSize,
		requests:    make(map[string][]time.Time),
	}
}

func (rl *slidingWindowRatelimiter) AllowRequest(key string) bool {
	currentTime := time.Now()
	queue, exists := rl.requests[key]
	if !exists {
		rl.requests[key] = make([]time.Time, 0)
		queue = rl.requests[key]
	}
	i := 0
	for i < len(queue) && currentTime.Sub(queue[i]) > rl.windowSize {
		i++
	}
	queue = queue[i:]
	if len(queue) < rl.maxRequests {
		queue = append(queue, currentTime)
		rl.requests[key] = queue
		return true
	}
	return false
}
