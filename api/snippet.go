package api

import (
	"blan-backend/database"
	"blan-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateSnippetRequest struct {
	Source string `json:"source" binding:"required"`
}

func CreateSnippetHandler(c *gin.Context) {
	var req CreateSnippetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userId lost"})
		return
	}

	snippet := models.Snippet{
		UserID: userID.(uint),
		Source: req.Source,
	}

	if err := database.DB.Create(&snippet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the code."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Code saved successfully", "snippet_id": snippet.ID})
}

func GetSnippetHandler(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "UserID lost."})
		return
	}

	var snippets []models.Snippet
	if err := database.DB.Where("user_id = ?", userID).Find(&snippets).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the snippets."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"count":    len(snippets),
		"snippets": snippets,
	})
}
