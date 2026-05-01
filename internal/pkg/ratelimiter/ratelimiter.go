package ratelimiter

import (
	"fmt"
	redis "learn/internal/config"
	"learn/internal/model"
	"learn/internal/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RateLimitWindow = 10 * time.Second
	RateLimitMaxReq = 100
)

// RateLimiterMiddleware checks and enforces a default rate limit per user/IP.
func RateLimiterMiddleware() gin.HandlerFunc {
	return Limit("default", RateLimitMaxReq, RateLimitWindow)
}

// Limit checks and enforces a route-specific rate limit per user/IP.
func Limit(name string, maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		identifier := identifierFromContext(c)
		rateLimitKey := fmt.Sprintf("ratelimit:%s:%s", name, identifier)

		count, err := redis.Rdb.Incr(redis.Ctx, rateLimitKey).Result()
		if err != nil {
			response.SendInternalServerError(c, nil, fmt.Errorf("rate limiter internal error: %w", err))
			return
		}

		if count == 1 {
			redis.Rdb.Expire(redis.Ctx, rateLimitKey, window)
		}

		remaining := maxRequests - count
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if count > maxRequests {
			c.Header("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
			response.SendTooManyRequestsError(c, fmt.Sprintf("Too many requests. Please try again after %s", window.String()))
			return
		}

		c.Next()
	}
}

func identifierFromContext(c *gin.Context) string {
	if userCtx, exists := c.Get("user"); exists {
		if authUser, ok := userCtx.(model.User); ok {
			return fmt.Sprintf("user:%d", authUser.ID)
		}
	}
	return "ip:" + c.ClientIP()
}
