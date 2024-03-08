package middleware

import (
	"github.com/andrelmm/goexpert-go-rate-limiter/pkg/limiter"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RateLimiterMiddleware(l *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the client is blocked. If it is, return a 429 status code
		if l.IsBlocked(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
			})
			c.Abort()
			return
		}

		// Check if the client has an API key. If it does, use it as the key for the rate limiter. If it doesn't, use the client's IP address
		apiKey := c.GetHeader("API_KEY")
		if apiKey != "" {
			if !l.CheckRateLimit(apiKey, true) {
				l.Block(apiKey)
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
				})
				c.Abort()
				return
			}
		} else {
			ip := c.ClientIP()
			if !l.CheckRateLimit(ip, false) {
				l.Block(ip)
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "you have reached the maximum number of requests or actions allowed within a certain time frame",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
