package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup) {

	auth := router.Group("/auth")

	{
		// Signup or register new user
		auth.POST("/signup", middlewares.RateLimiter("5-M"), handlers.Signup)

		// Signin and get token
		auth.POST("/signin", middlewares.RateLimiter("10-M"), handlers.Signin)

		// Logout the current user
		auth.POST("/logout", middlewares.RateLimiter("10-M"), handlers.Logout)
	}

}
