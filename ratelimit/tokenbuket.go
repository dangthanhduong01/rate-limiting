package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBucketRateLimiter struct {
	rdb      *redis.Client
	capacity int
	fillRate float64
}

func NewTokenBucketRateLimiter(rdb *redis.Client, capacity int, fillRate float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		rdb:      rdb,
		capacity: capacity,
		fillRate: fillRate,
	}
}
func (l *TokenBucketRateLimiter) Allow(userId string) bool {
	ctx := context.Background()

	tokenKey := "tb:tokens" + userId
	timeKey := "tb:timestamp" + userId

	now := float64(time.Now().UnixNano()) / 1e9

	tokenStr, _ := l.rdb.Get(ctx, tokenKey).Result()
	lastTimeStr, _ := l.rdb.Get(ctx, timeKey).Result()

	var tokens float64
	var lastTime float64

	if tokenStr != "" {
		tokens, _ = strconv.ParseFloat(tokenStr, 64)
	} else {
		tokens = float64(l.capacity)
	}

	if lastTimeStr != "" {
		lastTime, _ = strconv.ParseFloat(lastTimeStr, 64)
	} else {
		lastTime = now
	}

	elapsed := now - lastTime
	refilled := elapsed * l.fillRate
	tokens = min(tokens+refilled, float64(l.capacity))

	if tokens >= 1 {
		tokens -= 1

		l.rdb.Set(ctx, tokenKey, tokens, 0)
		l.rdb.Set(ctx, timeKey, now, 0)
		return true
	}
	// Not enough tokens
	return false
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
