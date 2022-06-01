package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// latencyLogger records the latency of set request handled
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

/*
func RequireAuthorization(l AuthorizationLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if !authorized(token, l) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// user is authorized, continue chaining HandlerFuncs
		c.Next()
	}
}

func authorize(token string, l AuthorizationLevel) bool {
	level := getLevelFromToken(token)
	return meetsPermissions(l, level)
}
*/
