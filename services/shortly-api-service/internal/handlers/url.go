package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"shortly-proto/gen/key"

	"shortly-api-service/internal/clients"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
	"shortly-api-service/internal/validators"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUrl(ctx *gin.Context) {

	idInterface, exists := ctx.Get("id")

	if !exists {
		utils.Log.Error("Id not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized: Id is missing from context",
		})
		return
	}

	id, ok := idInterface.(int)

	if !ok {
		utils.Log.Error("Failed to assert id type from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	idStr := strconv.Itoa(id)

	var data validators.CreateUrlValidator

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Log.Error("Failed to bind request body", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	validationErrors := validators.ValidateCreateUrlData(data)

	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":          false,
			"validation_error": validationErrors,
		})
		return
	}

	var existing models.Url

	if err := database.DB.Where("original_url = ?", data.OriginalURL).First(&existing).Error; err == nil {
		utils.Log.Error("Url is already shortened", "error", err)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "This URL has already been shortened.",
		})
		return
	}

	if data.ShortKey == "" {
		key, err := clients.KGSClient.GetKey(context.Background(), &key.Empty{})

		if err != nil {
			utils.Log.Error("Failed to get key from KGS service", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to generate short key",
			})
			return
		}

		data.ShortKey = key.Key
	}

	var existingKey models.Url

	if err := database.DB.Where("short_key = ?", data.ShortKey).First(&existingKey).Error; err == nil {
		utils.Log.Error("ShortKey already exists")
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "The generated ShortKey already exists, please try again",
		})
		return
	}

	newUrl := models.Url{
		OriginalURL: data.OriginalURL,
		ShortKey:    data.ShortKey,
		Title:       data.Title,
		UserID:      &idStr,
	}

	if err := database.DB.Create(&newUrl).Error; err != nil {
		utils.Log.Error("Failed to create URL", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create URL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"short_url": newUrl.ShortKey,
			"original":  newUrl.OriginalURL,
			"title":     newUrl.Title,
		},
		"message": "URL successfully created",
	})

}

func GetAllUrls(ctx *gin.Context) {

	idInterface, exists := ctx.Get("id")

	if !exists {
		utils.Log.Error("User ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized: User ID missing",
		})
		return
	}

	id, ok := idInterface.(int)

	if !ok {
		utils.Log.Error("Failed to assert user ID type")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	userID := strconv.Itoa(id)

	var urls []models.Url

	if err := database.DB.
		Where("user_id = ?", userID).
		Find(&urls).Error; err != nil {
		utils.Log.Error("Failed to fetch URLs", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve URLs",
		})
		return
	}

	response := make([]gin.H, 0, len(urls))

	for _, url := range urls {
		response = append(response, gin.H{
			"original_url": url.OriginalURL,
			"short_url":    url.ShortKey,
			"title":        url.Title,
			"clicks":       url.Clicks,
			"created_at":   url.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "URLs retrieved successfully",
	})
}

func GetUrlDetails(ctx *gin.Context) {

	shortKey := ctx.Param("shortKey")

	if shortKey == "" {
		utils.Log.Error("Short key is missing from path")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Short key is required",
		})
		return
	}

	var url models.Url

	if err := database.DB.Where("short_key = ?", shortKey).First(&url).Error; err != nil {
		utils.Log.Error("Failed to find URL", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "URL not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_url": url.OriginalURL,
			"short_url":    url.ShortKey,
			"title":        url.Title,
			"clicks":       url.Clicks,
			"created_at":   url.CreatedAt,
			"expires_at":   url.ExpiresAt,
		},
		"message": "URL details retrieved successfully",
	})
}

func RedirectToOriginalUrl(ctx *gin.Context) {

	shortKey := ctx.Param("shortKey")

	if shortKey == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Short key is required",
		})
		return
	}

	var url models.Url

	if err := database.DB.Where("short_key = ?", shortKey).First(&url).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	go func() {
		err := database.DB.Model(&url).UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error
		if err != nil {
			utils.Log.Error("Failed to update click count", "error", err)
		}
	}()

	go func() {
		analytics := models.Analytics{
			UrlID:     strconv.FormatUint(uint64(url.ID), 10),
			ClickedAt: time.Now(),
			IPAddress: ctx.ClientIP(),
			UserAgent: ctx.GetHeader("User-Agent"), // Capture user-agent
			Referrer:  ctx.GetHeader("Referer"),    // Capture referer
			Country:   "",                          // Optional: Implement GeoIP or similar to get country from IP
			Device:    "",                          // Optional: Use user-agent parser to detect device
			Browser:   "",                          // Optional: Use user-agent parser to detect browser
			OS:        "",                          // Optional: Use user-agent parser to detect OS
		}

		if err := database.DB.Create(&analytics).Error; err != nil {
			utils.Log.Error("Failed to store analytics", "error", err)
		}
	}()

	ctx.Redirect(http.StatusFound, url.OriginalURL)
}

func UpdateUrl(ctx *gin.Context) {

}

func DeleteUrl(ctx *gin.Context) {

}
