package handlers

import (
	"net/http"
	"strings"

	"shortly-api-service/internal/database"
	"shortly-api-service/internal/dto"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
	"shortly-api-service/internal/validators"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Signup(ctx *gin.Context) {

	var data validators.SignupValidator

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Log.Error("Failed to bind request body", "error", err)
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
		utils.Log.Error("Validation failed on signup", "error", validationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":          false,
			"validation_error": validationErrors,
		})
		return
	}

	var existingUser models.User

	if err := database.DB.Where("email = ?", data.Email).First(&existingUser).Error; err == nil {
		utils.Log.Warn("Signup attempt with existing email", "email", data.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "User already exists",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.Log.Error("Database error", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database error",
		})
		return
	}

	hashPassword, err := utils.HashPassword(data.Password)

	if err != nil {
		utils.Log.Error("Error hashing password", "error", err)
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
		utils.Log.Error("Failed to create user", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create user",
		})
		return
	}

	utils.Log.Info("User signed up successfully", "user_id", user.ID, "email", user.Email)

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			Username:  data.Username,
			CreatedAt: user.CreatedAt,
		},
		"message": "User registered successfully",
	})

}

func Signin(ctx *gin.Context) {

	var data validators.SigninValidator

	if err := ctx.ShouldBindJSON(&data); err != nil {
		utils.Log.Error("Failed to bind request body", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	data.Email = strings.TrimSpace(strings.ToLower(data.Email))

	validationErrors := validators.ValidateSigninData(data)

	if len(validationErrors) > 0 {
		utils.Log.Error("Validation failed on signin", "error", validationErrors)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success":          false,
			"validation_error": validationErrors,
		})
		return
	}

	var user models.User

	if err := database.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		utils.Log.Warn("Signin attempt with unregistered email", "email", data.Email)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	isValid := utils.VerifyPassword(data.Password, user.Password)

	if !isValid {
		utils.Log.Error("Password is not valid")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Password is not valid",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)

	if err != nil {
		utils.Log.Error("Could not generate token")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Could not generate token",
		})
		return
	}

	ctx.SetCookie("token", token, 86400, "/", "", true, true)

	utils.Log.Info("User login attempt",
		"email", data.Email,
		"ip", ctx.ClientIP(),
		"user_agent", ctx.Request.UserAgent(),
	)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
		"data": dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		"message": "Login successful",
	})

}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "", true, true)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successfully",
	})
}
