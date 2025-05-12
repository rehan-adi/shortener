package routes

import (
	"shortly-api-service/internal/handlers"
	"shortly-api-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func UrlRouter(router *gin.RouterGroup) {

	// Redirect to Original Url
	router.GET("/url/redirect/:shortKey", middlewares.RateLimiter("50-m"), handlers.RedirectToOriginalUrl)

	url := router.Group("/url").Use(middlewares.AuthMiddleware())

	{
		// Get all URLs (for login user)
		url.GET("/", middlewares.RateLimiter("20-M"), handlers.GetAllUrls)

		// Shorten a URL
		url.POST("/shorten", middlewares.RateLimiter("5-M"), handlers.CreateUrl)

		// Get URL details by shortKey
		url.GET("/:shortKey", middlewares.RateLimiter("20-M"), handlers.GetUrlDetails)

		// Update an existing URL
		url.PATCH("/:shortKey", middlewares.RateLimiter("5-M"), handlers.UpdateUrl)

		// Delete a URL
		url.DELETE("/:shortKey", middlewares.RateLimiter("5-M"), handlers.DeleteUrl)
	}

}
