package ratelimiter

import (
	"fmt"
	"learn/internal/feature/auth"
	redis "learn/internal/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RateLimitWindow = 10 * time.Second
	RateLimitMaxReq = 100
)

// RateLimiterMiddleware checks and enforces rate limits per user/IP.
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var identifier string

		// Try to get user from context (if authenticated)
		user, exists := c.Get("user")
		if exists {
			// Assuming user is of type auth.User and has a Name field
			if authUser, ok := user.(auth.User); ok {
				identifier = authUser.Name
			} else {
				// Fallback to IP if user object is not as expected
				identifier = c.ClientIP()
			}
		} else {
			// For unauthenticated routes, use client IP
			identifier = c.ClientIP()
		}

		rateLimitKey := fmt.Sprintf("ratelimit:%s", identifier)

		// Increment counter for the identifier
		count, err := redis.Rdb.Incr(redis.Ctx, rateLimitKey).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter internal error"})
			c.Abort()
			return
		}

		// If it's the first request in this window, set expiry
		if count == 1 {
			redis.Rdb.Expire(redis.Ctx, rateLimitKey, RateLimitWindow)
		}

		// Check if the limit is exceeded
		if count > RateLimitMaxReq {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": fmt.Sprintf("Too many requests. Please try again after %s", RateLimitWindow.String())})
			c.Abort()
			return
		}

		c.Next()
	}
}