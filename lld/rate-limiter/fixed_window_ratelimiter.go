package rate_limiter

import (
	"sync"
	"time"
)

type fixedWindowRatelimiter struct {
	maxRequests int
	windowSize  time.Duration
	// make map sync
	requests sync.Map
	mu       sync.Mutex
}

type requestData struct {
	requestCount    int
	windowStartTime time.Time
}

func NewFixedWindowRatelimiter(options ...func(*fixedWindowRatelimiter)) *fixedWindowRatelimiter {
	ratelimiter := &fixedWindowRatelimiter{
		maxRequests: 0,
		windowSize:  0,
		requests:    sync.Map{},
	}
	for _, o := range options {
		o(ratelimiter)
	}
	return ratelimiter
}

func WithMaxRequests(maxRequests int) func(*fixedWindowRatelimiter) {
	return func(ratelimiter *fixedWindowRatelimiter) {
		ratelimiter.maxRequests = maxRequests
	}
}
func WithWindowSize(windowSize time.Duration) func(*fixedWindowRatelimiter) {
	return func(ratelimiter *fixedWindowRatelimiter) {
		ratelimiter.windowSize = windowSize
	}
}

func (ratelimiter *fixedWindowRatelimiter) AllowRequest(key string) bool {
	currentTime := time.Now()
	reqDataAny, exists := ratelimiter.requests.Load(key)
	reqData, _ := reqDataAny.(requestData)
	if !exists {
		reqData.requestCount = 1
		reqData.windowStartTime = currentTime
		ratelimiter.requests.Store(key, reqData)
		return true
	}
	if currentTime.Sub(reqData.windowStartTime) > ratelimiter.windowSize {
		reqData.requestCount = 1
		reqData.windowStartTime = currentTime
		ratelimiter.requests.Store(key, reqData)
		return true
	}
	if reqData.requestCount < ratelimiter.maxRequests {
		reqData.requestCount++
		ratelimiter.requests.Store(key, reqData)
		return true
	}
	return false
}
