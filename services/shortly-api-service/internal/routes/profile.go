package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func ProfileRouter(router *gin.RouterGroup) {

	profile := router.Group("/profile")

	{
		profile.GET("/", middlewares.AuthMiddleware(), handlers.GetUserProfile)
		profile.PATCH("/update", middlewares.AuthMiddleware(), handlers.UpdateUserProfile)
	}

}
