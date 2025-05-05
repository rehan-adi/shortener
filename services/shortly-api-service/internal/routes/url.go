package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UrlRouter(router *gin.RouterGroup) {

	url := router.Group("/url").Use(middlewares.AuthMiddleware())

	{
		// Get all URLs (for login user)
		url.GET("/", handlers.GetAllUrls)

		// Shorten a URL
		url.POST("/shorten", handlers.CreateUrl)

		// Get URL details by shortKey
		url.GET("/:shortKey", handlers.GetUrlDetails)

		// Update an existing URL
		url.PATCH("/:shortKey", handlers.UpdateUrl)

		// Delete a URL
		url.DELETE("/:shortKey", handlers.DeleteUrl)
	}

	url.GET("/redirect/:shortKey", handlers.RedirectToOriginalUrl)

}
