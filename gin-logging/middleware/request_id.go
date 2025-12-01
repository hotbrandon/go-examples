package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Why check X-Request-ID even though RequestID() runs first?
// ref: readme.md
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("request_id", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}
