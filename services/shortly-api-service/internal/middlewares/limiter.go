package middlewares

import (
	"context"
	"shortly-api-service/config"
	"shortly-api-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	redisStore "github.com/ulule/limiter/v3/drivers/store/redis"

	"github.com/redis/go-redis/v9"
)

func RateLimiter(rateString string) gin.HandlerFunc {

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.REDIS_ADDR,
		Password: "",
		DB:       2,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic("🔴 Failed to connect to Redis for rate limiting: " + err.Error())
	}

	rate, err := limiter.NewRateFromFormatted(rateString)

	if err != nil {
		panic("🔴 Invalid rate limit format: " + err.Error())
	}

	store, err := redisStore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})

	if err != nil {
		panic("🔴 Failed to create rate limit store: " + err.Error())
	}

	instance := limiter.New(store, rate)

	return func(c *gin.Context) {

		contextKey := c.ClientIP()

		ctx, err := instance.Get(c, contextKey)

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Rate limiter internal error"})
			return
		}

		if ctx.Reached {
			utils.Log.Warn("Rate limit exceeded", "ip", contextKey, "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(429, gin.H{"error": "Too Many Requests"})
			return
		}

		c.Next()
	}
}
