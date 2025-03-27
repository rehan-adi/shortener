package routes

import (
	"shortly-api-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func HealthRouter(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		health.GET("/", handlers.HealthCheck)
	}
}