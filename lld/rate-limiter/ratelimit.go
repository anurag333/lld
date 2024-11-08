package rate_limiter

type RateLimiter interface {
	AllowRequest(string) bool
}
