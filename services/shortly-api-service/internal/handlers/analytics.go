package handlers

import (
	"net/http"
	"shortly-api-service/internal/database"
	"shortly-api-service/internal/dto"
	"shortly-api-service/internal/models"
	"shortly-api-service/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAnalytics(ctx *gin.Context) {

	idInterface, exists := ctx.Get("id")

	if !exists {
		utils.Log.Error("Id not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized: Id is missing from context",
		})
		return
	}

	urlId := ctx.Param("urlId")

	if urlId == "" {
		utils.Log.Error("Missing urlId in request path")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing urlId in path",
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

	if err := database.DB.Where("id = ?", idStr).Error; err != nil {
		utils.Log.Error("User not found", "error", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
	}

	var analytics []models.Analytics

	if err := database.DB.Where("url_id  = ?", urlId).Order("clicked_at desc").Find(&analytics).Error; err != nil {
		utils.Log.Error("Failed to fetch analytics", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch analytics",
		})
	}

	var response []dto.AnalyticsResponse = make([]dto.AnalyticsResponse, 0)

	for _, a := range analytics {
		response = append(response, dto.AnalyticsResponse{
			IPAddress: a.IPAddress,
			OS:        a.OS,
			Device:    a.Device,
			Browser:   a.Browser,
			UserAgent: a.UserAgent,
			Referrer:  a.Referrer,
			Country:   a.Country,
			ClickedAt: a.ClickedAt.Format("2006-01-02 15:04:05"),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "Analytics retrieved successfully",
	})

}
