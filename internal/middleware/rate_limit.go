package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter is a fixed‑window counter (cheap and good enough for APIs).
// For production‑grade accuracy you might swap to a token‑bucket Lua script.
func RateLimiter(rdb *redis.Client, window time.Duration, max int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)
		sec := time.Now().UTC().Truncate(window).Unix()
		key := fmt.Sprintf("rl:%d:%s", sec, ip)

		// INCR and set TTL at the same time (pipeline = 1 RTT)
		pipe := rdb.TxPipeline()
		cnt := pipe.Incr(c, key)
		pipe.Expire(c, key, window+time.Second) // keep key slightly longer
		_, _ = pipe.Exec(c)

		if cnt.Val() > int64(max) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprint(max))
		c.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprint(max-int(cnt.Val())))
		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	if ip := c.ClientIP(); ip != "" {
		// strip port if gin returned host:port (unix sockets)
		if host, _, err := net.SplitHostPort(ip); err == nil {
			return host
		}
		return ip
	}
	return "unknown"
}
