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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user identity missing from context"})
		return
	}

	snippet := models.Snippet{
		UserID: userID.(uint),
		Source: req.Source,
	}

	if err := database.DB.Create(&snippet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save snippet"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "snippet saved successfully", "snippet_id": snippet.ID})
}

func GetSnippetHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user identity missing from context"})
		return
	}

	var snippets []models.Snippet
	if err := database.DB.Where("user_id = ?", userID).Find(&snippets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch snippets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":    len(snippets),
		"snippets": snippets,
	})
}
