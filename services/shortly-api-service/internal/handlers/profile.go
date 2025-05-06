package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"shortly-api-service/internal/database"
	"shortly-api-service/internal/dto"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/redis"
	"shortly-api-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(ctx *gin.Context) {

	emailInterface, exists := ctx.Get("email")

	if !exists {
		utils.Log.Error("Email not found")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized: Email missing",
		})
		return
	}

	email, ok := emailInterface.(string)

	if !ok {
		utils.Log.Error("Failed to assert email type from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	cacheKey := "user:profile:" + email

	cached, err := redis.RedisClient.Get(ctx.Request.Context(), cacheKey).Result()

	if err == nil && cached != "" {
		var cachedDTO dto.UserDTO
		if err := json.Unmarshal([]byte(cached), &cachedDTO); err == nil {
			utils.Log.Info("User profile served from Redis")
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    cachedDTO,
				"message": "Profile retrieved successfully",
			})
			return
		}
	}

	var user models.User

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.Log.Error("No user found", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	jsonBytes, err := json.Marshal(userDTO)

	if err != nil {
		utils.Log.Error("Failed to marshal user DTO", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	if err := redis.RedisClient.Set(
		ctx.Request.Context(),
		cacheKey,
		jsonBytes,
		30*time.Minute).Err(); err != nil {
		utils.Log.Error("Failed to cache user profile", "error", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userDTO,
		"message": "Profile retrieved successfully",
	})

}

func UpdateUserProfile(ctx *gin.Context) {

	emailInterface, exists := ctx.Get("email")

	if !exists {
		utils.Log.Error("Email not found")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized: Email missing",
		})
		return
	}

	email, ok := emailInterface.(string)

	if !ok {
		utils.Log.Error("Failed to assert email type from context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	var data dto.UpdateUserDTO

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Log.Warn("Invalid input for user update", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input: " + err.Error(),
		})
		return
	}

	var user models.User

	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.Log.Error("User not found during update", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	user.Username = strings.TrimSpace(data.Username)

	if user.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Username cannot be empty",
		})
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.Log.Error("Failed to update user profile", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update profile",
		})
		return
	}

	cacheKey := "user:profile:" + email

	userDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	jsonBytes, err := json.Marshal(userDTO)

	if err != nil {
		utils.Log.Error("Failed to marshal user DTO", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	if err := redis.RedisClient.Set(
		ctx.Request.Context(),
		cacheKey,
		jsonBytes,
		30*time.Minute).Err(); err != nil {
		utils.Log.Error("Failed to update Redis cache", "error", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userDTO,
		"message": "Profile updated successfully",
	})

}
