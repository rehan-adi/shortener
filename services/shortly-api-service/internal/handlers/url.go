package handlers

import (
	"context"
	"net/http"
	"strconv"
	
	"shortly-proto/gen/key"

	"shortly-api-service/internal/clients"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
	"shortly-api-service/internal/validators"

	"github.com/gin-gonic/gin"
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

}

func GetUrlDetails(ctx *gin.Context) {

}

func UpdateUrl(ctx *gin.Context) {

}

func DeleteUrl(ctx *gin.Context) {

}
