package routes

import (
	"shortly-api-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup) {

	auth := router.Group("/auth")

	{
		// Signup or register new user
		auth.POST("/signup", handlers.Signup)

		// Signin and get token
		auth.POST("/signin", handlers.Signin)

		// Logout the current user
		auth.POST("/logout", handlers.Logout)
	}

}
