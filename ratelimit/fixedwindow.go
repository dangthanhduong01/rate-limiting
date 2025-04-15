package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type FixedWindowRateLimiter struct {
	rdb       *redis.Client
	limit     int
	windowSec int
}

func NewFixedWindowRateLimiter(rdb *redis.Client, limit, windowSec int) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{rdb, limit, windowSec}
}

func (l *FixedWindowRateLimiter) Allow(userId string) bool {
	ctx := context.Background()
	now := time.Now().Unix()

	window := now - (now % int64(l.windowSec))
	key := fmt.Sprintf("fw:%s:%s", userId, string(rune(window)))

	count, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false
	}
	if count == 1 {
		l.rdb.Expire(ctx, key, time.Duration(l.windowSec)*time.Second)
	}
	if count > int64(l.limit) {
		return false
	}
	return true
}
