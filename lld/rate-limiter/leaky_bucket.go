package rate_limiter

import (
	"sync"
	"time"
)

// leakyBucket represents a leaky bucket rate limiter for a specific customer
type leakyBucket struct {
	capacity   int           // Maximum capacity of the bucket
	rps        int           // Requests per second allowed to pass
	queue      chan struct{} // Channel acting as a queue for incoming requests
	mu         sync.Mutex    // Mutex for safe concurrent access
	stopSignal chan struct{} // Channel to signal when to stop leaking
}

// newLeakyBucket initializes a new leakyBucket with given capacity and rps
func newLeakyBucket(capacity, rps int) *leakyBucket {
	bucket := &leakyBucket{
		capacity:   capacity,
		rps:        rps,
		queue:      make(chan struct{}, capacity),
		stopSignal: make(chan struct{}),
	}
	go bucket.startLeaking()
	return bucket
}

// allowRequest checks if a request is allowed based on the rate limit
func (b *leakyBucket) allowRequest() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	select {
	case b.queue <- struct{}{}: // Add a request to the queue if capacity allows
		return true
	default:
		return false // Bucket is full, request is not allowed
	}
}

// startLeaking continuously drains the bucket at the defined rate (rps)
func (b *leakyBucket) startLeaking() {
	ticker := time.NewTicker(time.Second / time.Duration(b.rps))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.mu.Lock()
			select {
			case <-b.queue: // Remove a request from the queue
			default:
				// No requests in the queue to process
			}
			b.mu.Unlock()
		case <-b.stopSignal:
			return // stop leaking if we receive a stop signal
		}
	}
}

// stop stops the leaky bucket from leaking
func (b *leakyBucket) stop() {
	close(b.stopSignal)
}

// LeakyBucketRatelimiter manages rate limiting for multiple customers, each with their own leakyBucket
type LeakyBucketRatelimiter struct {
	buckets   map[string]*leakyBucket // Map of customer ID to their respective leakyBucket
	bucketCap int                     // Bucket capacity for each customer
	bucketRPS int                     // Requests per second for each customer's bucket
	mu        sync.Mutex              // Mutex for safe access to the buckets map
}

// NewRateLimiter initializes a new LeakyBucketRatelimiter
func NewLeakyBucketRatelimiter(bucketCap, bucketRPS int) *LeakyBucketRatelimiter {
	return &LeakyBucketRatelimiter{
		buckets:   make(map[string]*leakyBucket),
		bucketCap: bucketCap,
		bucketRPS: bucketRPS,
	}
}

// getOrCreateBucket retrieves or creates a leakyBucket for a given customer ID
func (rl *LeakyBucketRatelimiter) getOrCreateBucket(customerID string) *leakyBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if a bucket exists for the customer; if not, create one
	bucket, exists := rl.buckets[customerID]
	if !exists {
		bucket = newLeakyBucket(rl.bucketCap, rl.bucketRPS)
		rl.buckets[customerID] = bucket
	}
	return bucket
}

// AllowRequest checks if a request is allowed for the given customer ID
func (rl *LeakyBucketRatelimiter) AllowRequest(customerID string) bool {
	bucket := rl.getOrCreateBucket(customerID)
	return bucket.allowRequest()
}

// StopAll stops all leaky buckets managed by the rate limiter
func (rl *LeakyBucketRatelimiter) StopAll() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for _, bucket := range rl.buckets {
		bucket.stop()
	}
}
