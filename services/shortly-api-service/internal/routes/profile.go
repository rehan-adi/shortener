package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func ProfileRouter(router *gin.RouterGroup) {

	profile := router.Group("/profile").Use(middlewares.AuthMiddleware())

	{
		// Get the authenticated user's profile information
		profile.GET("/", middlewares.RateLimiter("20-M"), handlers.GetUserProfile)

		// Update the authenticated user's profile information
		profile.PATCH("/update", middlewares.RateLimiter("5-M"), handlers.UpdateUserProfile)
	}

}
