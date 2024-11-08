package rate_limiter

import (
	"math"
	"time"
)

type tokenBucketRatelimiter struct {
	bucketSize  int
	refillEvery time.Duration
	requests    map[string]tokenData
}

type tokenData struct {
	freeTokens           int
	lastFillingTimestamp time.Time
}

func NewTokenBucketRatelimiter(bucketSize int, refillEvery time.Duration) *tokenBucketRatelimiter {
	return &tokenBucketRatelimiter{
		bucketSize:  bucketSize,
		refillEvery: refillEvery,
		requests:    make(map[string]tokenData),
	}
}

func (rl *tokenBucketRatelimiter) AllowRequest(key string) bool {
	currentTime := time.Now()
	token, exists := rl.requests[key]
	if !exists {
		rl.requests[key] = tokenData{
			freeTokens:           rl.bucketSize - 1,
			lastFillingTimestamp: currentTime,
		}
		return true
	}
	newTokens := float64((currentTime.Sub(token.lastFillingTimestamp))/rl.refillEvery) * float64(rl.bucketSize)
	if newTokens > 0 {
		token.lastFillingTimestamp = currentTime
		token.freeTokens = int(math.Max(float64(token.freeTokens)+newTokens, float64(rl.bucketSize)))
	}
	if token.freeTokens > 0 {
		token.freeTokens--
		rl.requests[key] = token
		return true
	}
	return false
}
