package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AnalyticsRouter(router *gin.RouterGroup) {

	analytics := router.Group("/analytics").Use(middlewares.AuthMiddleware())

	{
		analytics.GET("/:urlId", middlewares.RateLimiter("10-m"), handlers.GetAnalytics)
	}

}
