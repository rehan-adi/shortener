package routes

import (
	"shortly-api-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.RouterGroup) {

	auth := router.Group("/auth")

	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/signin", handlers.Signin)
		auth.POST("/logout", handlers.Logout)
	}

}