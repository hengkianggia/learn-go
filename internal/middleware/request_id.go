package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const requestIDKey = "request_id"

// RequestIDMiddleware propagates or creates a request ID for tracing logs and responses.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}

		c.Set(requestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		ctx := context.WithValue(c.Request.Context(), requestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if value, ok := c.Get(requestIDKey); ok {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
