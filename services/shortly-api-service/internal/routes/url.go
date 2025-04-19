package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UrlRouter(router *gin.RouterGroup) {

	url := router.Group("/url").Use(middlewares.AuthMiddleware())

	{
		url.GET("/", handlers.GetAllUrls)
		url.POST("/shorten", handlers.CreateUrl)
		url.GET("/:shortKey", handlers.GetUrlDetails)
		url.PATCH("/:shortKey", handlers.UpdateUrl)
		url.DELETE("/:shortKey", handlers.DeleteUrl)
	}
}
