package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LatencyLogger records the latency of set request handled
func LatencyLogger(l *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()

		c.Next()

		// after request
		latency := time.Since(t)
		l.Printf("%v request on %q returned status %v in %v", c.Request.Method, c.Request.URL, c.Writer.Status(), latency)
	}
}
