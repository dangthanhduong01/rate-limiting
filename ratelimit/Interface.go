package ratelimit

type RateLimiter interface {
	Allow(userId string) bool
}
