package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"shortly-proto/gen/key"

	"shortly-api-service/internal/clients"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/dto"
	"shortly-api-service/internal/lib"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/redis"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
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

	if err := database.DB.Where("original_url = ?  AND user_id = ? ", data.OriginalURL, idStr).First(&existing).Error; err == nil {
		utils.Log.Error("Url is already shortened")
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

	cacheKey := "url:" + data.ShortKey

	jsonByte, err := json.Marshal(newUrl)

	if err != nil {
		utils.Log.Error("Failed to marshal url data")
	}

	go func() {
		err := redis.RedisClient.Set(
			ctx.Request.Context(),
			cacheKey,
			jsonByte,
			24*time.Hour).Err()

		if err != nil {
			utils.Log.Error("Failed to cache generated url data", "error", err)
		}
	}()

	utils.Log.Info("URL successfully created", "shortKey", data.ShortKey, "userID", idStr)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.CreateUrlResponseDTO{
			ID:          newUrl.ID,
			OriginalURL: newUrl.OriginalURL,
			ShortKey:    newUrl.ShortKey,
			Title:       newUrl.Title,
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

	response := make([]dto.GetUrlResponseDTO, 0, len(urls))

	for _, url := range urls {
		response = append(response, dto.GetUrlResponseDTO{
			ID:          url.ID,
			OriginalURL: url.OriginalURL,
			ShortKey:    url.ShortKey,
			Title:       url.Title,
			Clicks:      url.Clicks,
			CreatedAt:   url.CreatedAt,
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
		"data": dto.GetUrlResponseDTO{
			ID:          url.ID,
			OriginalURL: url.OriginalURL,
			ShortKey:    url.ShortKey,
			Title:       url.Title,
			Clicks:      url.Clicks,
			CreatedAt:   url.CreatedAt,
		},
		"message": "URL details retrieved successfully",
	})
}

func RedirectToOriginalUrl(ctx *gin.Context) {

	shortKey := ctx.Param("shortKey")

	if shortKey == "" {
		utils.Log.Error("Short key is missing from path")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Short key is required",
		})
		return
	}

	cacheKey := "url:" + shortKey

	data, err := redis.RedisClient.Get(ctx.Request.Context(), cacheKey).Result()

	if err == nil && data != "" {
		var cachedDTO models.Url
		if err := json.Unmarshal([]byte(data), &cachedDTO); err == nil {
			utils.Log.Info("URL served from Redis cache")

			go incrementClickCount(cachedDTO.ID)
			go storeAnalytics(ctx, cachedDTO.ID)

			ctx.Redirect(http.StatusFound, cachedDTO.OriginalURL)
			return
		}
	}

	var url models.Url

	if err := database.DB.Where("short_key = ?", shortKey).First(&url).Error; err != nil {
		utils.Log.Error("Short Key not found in database", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "URL not found",
		})
		return
	}

	go func() {
		jsonData, err := json.Marshal(url)
		if err == nil {
			err = redis.RedisClient.Set(ctx.Request.Context(), cacheKey, jsonData, 24*time.Hour).Err()
			if err != nil {
				utils.Log.Error("Failed to cache URL in Redis", "error", err)
			}
		}
	}()

	// Async operations
	go incrementClickCount(url.ID)
	go storeAnalytics(ctx, url.ID)

	ctx.Redirect(http.StatusFound, url.OriginalURL)
}

func incrementClickCount(urlId uint) {
	err := database.DB.Model(&models.Url{}).Where("id = ?", urlId).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error
	if err != nil {
		utils.Log.Error("Failed to update click count", "error", err)
	}
}

func storeAnalytics(ctx *gin.Context, urlID uint) {
	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	country := lib.GetCountryFromIP(ip)
	device, browser, os := lib.ParseUserAgent(userAgent)

	analytics := models.Analytics{
		UrlID:     strconv.FormatUint(uint64(urlID), 10),
		ClickedAt: time.Now(),
		IPAddress: ip,
		UserAgent: userAgent,
		Referrer:  ctx.GetHeader("Referer"),
		Country:   country,
		Device:    device,
		Browser:   browser,
		OS:        os,
	}

	if err := database.DB.Create(&analytics).Error; err != nil {
		utils.Log.Error("Failed to store analytics", "error", err)
	}
}

func UpdateUrl(ctx *gin.Context) {

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
		utils.Log.Error("URL not found", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "URL not found",
		})
		return
	}

	var updateData validators.UpdateUrlValidator

	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		utils.Log.Error("Invalid input data", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input data",
		})
		return
	}
	validationErrors := validators.ValidateUpdateUrlData(updateData)

	if len(validationErrors) > 0 {
		utils.Log.Error("Validation errors", "errors", validationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   validationErrors,
		})
		return
	}

	if updateData.ShortKey != "" && updateData.ShortKey != shortKey {
		var existing models.Url
		if err := database.DB.Where("short_key = ?", updateData.ShortKey).First(&existing).Error; err == nil {
			utils.Log.Error("Short key already exists", "short_key", updateData.ShortKey)
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"error":   "This short key already exists.",
			})
			return
		}
	}

	if updateData.ShortKey != "" {
		url.ShortKey = updateData.ShortKey
	}
	if updateData.Title != "" {
		url.Title = updateData.Title
	}

	if err := database.DB.Save(&url).Error; err != nil {
		utils.Log.Error("Failed to update URL", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update URL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"original_url": url.OriginalURL,
			"short_url":    url.ShortKey,
			"title":        url.Title,
		},
		"message": "URL updated successfully",
	})
}

func DeleteUrl(ctx *gin.Context) {

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
		utils.Log.Error("URL not found", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "URL not found",
		})
		return
	}

	if err := database.DB.Delete(&url).Error; err != nil {
		utils.Log.Error("Failed to delete URL", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete URL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "URL deleted successfully",
	})
}
