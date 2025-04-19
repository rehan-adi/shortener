package middlewares

import (
	"net/http"
	"strings"

	"shortly-api-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		var err error

		token, err = ctx.Cookie("token")
		if err != nil {
			// If cookie is missing, check the Authorization header
			authHeader := ctx.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		// If token is still empty, return unauthorized
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized: No token provided"})
			ctx.Abort()
			return
		}

		// Verify the token
		claims, err := utils.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized: Invalid token"})
			ctx.Abort()
			return
		}

		// Extract user_id and email from claims
		userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized: Invalid user ID"})
			ctx.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized: Invalid email"})
			ctx.Abort()
			return
		}

		// Store in context for later use in routes
		ctx.Set("id", int(userID)) // Convert float64 to int
		ctx.Set("email", email)

		// Continue request processing
		ctx.Next()
	}
}
