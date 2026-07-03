package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "requestId"

// Logger returns a Gin middleware that logs each request with structured JSON
// via slog, including a unique request ID.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.New().String()
		c.Set(RequestIDKey, reqID)
		c.Header("X-Request-ID", reqID)

		start := time.Now()

		c.Next()

		slog.Info("request",
			slog.String("requestId", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("clientIP", c.ClientIP()),
		)
	}
}
