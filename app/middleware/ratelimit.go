package middleware

import (
	"net/http"
	"strconv"

	"github.com/dariubs/scaffold/app/utils"
	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimit creates a rate limiting middleware
func RateLimit(rate string) gin.HandlerFunc {
	rateLimit, err := limiter.NewRateFromFormatted(rate)
	if err != nil {
		utils.Logger.Error("Failed to parse rate limit", "err", err)
		return func(c *gin.Context) {
			c.Next()
		}
	}

	store := memory.NewStore()
	instance := limiter.New(store, rateLimit)

	return func(c *gin.Context) {
		key := c.ClientIP()

		context, err := instance.Get(c, key)
		if err != nil {
			utils.Logger.Error("Rate limiter error", "err", err)
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
