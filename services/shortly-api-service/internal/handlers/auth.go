package handlers

import (
	"net/http"
	"strings"

	"shortly-api-service/internal/database"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
	"shortly-api-service/internal/validators"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Signup(ctx *gin.Context) {

	var data validators.SignupValidator

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Log.Errorf("Failed to bind request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	data.Email = strings.TrimSpace(strings.ToLower(data.Email))
	data.Username = strings.TrimSpace(data.Username)

	validationErrors := validators.ValidateSignupData(data)

	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":          false,
			"validation_error": validationErrors,
		})
		return
	}

	var existingUser models.User

	if err := database.DB.Where("email = ?", data.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "User already exists",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.Log.Errorf("Database error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	hashPassword, err := utils.HashPassword(data.Password)

	if err != nil {
		utils.Log.Errorf("Error hashing password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Error hashing password",
		})
		return
	}

	user := models.User{
		Email:    data.Email,
		Username: data.Username,
		Password: hashPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Log.Errorf("Failed to create user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
	})

}