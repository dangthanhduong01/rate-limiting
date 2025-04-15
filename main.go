package main

import (
	"context"
	"net/http"
	"ratelimiter/ratelimit"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(limiter ratelimit.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Request.Header.Get("X-User-ID")

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
			c.Abort()
			return
		}

		if !limiter.Allow(userId) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	router := gin.Default()

	// Choose a rate limiter
	fixedWindowLimiter := ratelimit.NewFixedWindowRateLimiter(rdb, 5, 60)
	slidingWindowLimiter := ratelimit.NewSlidingWindowRateLimiter(rdb, 10, 60, 10)
	tokenBucketlimiter := ratelimit.NewTokenBucketRateLimiter(rdb, 10, 1)

	router.GET("/fixed", RateLimitMiddleware(fixedWindowLimiter), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "fixed window rate limit passed"})
	})

	router.GET("/sliding", RateLimitMiddleware(slidingWindowLimiter), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "sliding window rate limit passed"})
	})
	router.GET("/token", RateLimitMiddleware(tokenBucketlimiter), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "token bucket rate limit passed"})
	})
	router.Run(":8080")

}
