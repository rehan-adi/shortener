package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UrlRouter(router *gin.RouterGroup) {

	url := router.Group("/url").Use(middlewares.AuthMiddleware())

	{
		url.GET("/", middlewares.AuthMiddleware(), handlers.GetAllUrls)
		url.POST("/shorten", middlewares.AuthMiddleware(), handlers.CreateUrl)
		url.GET("/:shortKey", middlewares.AuthMiddleware(), handlers.GetUrlDetails)
		url.PATCH("/:shortKey", middlewares.AuthMiddleware(), handlers.UpdateUrl)
		url.DELETE("/:shortKey", middlewares.AuthMiddleware(), handlers.DeleteUrl)
	}
}
