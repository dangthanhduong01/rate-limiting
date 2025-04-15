package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SlidingWindowRateLimiter struct {
	rdb        *redis.Client
	limit      int
	windowSize int
	slotSize   int
}

func NewSlidingWindowRateLimiter(rdb *redis.Client, limit, windowSize, slotSize int) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{rdb, limit, windowSize, slotSize}
}
func (l *SlidingWindowRateLimiter) Allow(userId string) bool {
	ctx := context.Background()
	now := time.Now()
	slots := l.windowSize / l.slotSize

	var sum int64
	for i := 0; i < slots; i++ {
		slotTime := now.Add(time.Duration(-int64(i*l.slotSize)) * time.Second)
		slotKey := fmt.Sprintf("sw:%s:%d", userId, slotTime.Unix()/int64(l.slotSize))

		val, _ := l.rdb.Get(ctx, slotKey).Int64()
		sum += val
	}

	if sum >= int64(l.limit) {
		return false
	}

	currSlotKey := fmt.Sprintf("sw:%s:%d", userId, now.Unix()/int64(l.slotSize))
	l.rdb.Incr(ctx, currSlotKey)
	l.rdb.Expire(ctx, currSlotKey, time.Duration(l.windowSize)*time.Second)

	return true
}
