package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Extract request ID from context (provided by RequestID middleware)
		reqID, _ := c.Get("request_id")

		// Create a request-scoped logger
		log := base.With(
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_ip", c.ClientIP(),
		)

		// Make logger available to handlers
		c.Set("logger", log)

		c.Next()

		// After handler runs, log completion
		log.Info("request completed",
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
