package handlers

import (
	"net/http"
	"strings"

	"shortly-api-service/internal/database"
	"shortly-api-service/internal/dto"
	"shortly-api-service/internal/models"
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

	var user models.User

	if err := database.DB.Preload("Urls").Where("email = ?", email).First(&user).Error; err != nil {
		utils.Log.Error("No user found", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	UserDTO := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		UrlsCount: len(user.Urls),
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    UserDTO,
		"message": "Profile retrive successfully",
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

	if err := database.DB.Save(&user).Error; err != nil {
		utils.Log.Error("Failed to update user profile", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
	})

}
