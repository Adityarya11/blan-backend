package api

import (
	"blan-backend/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StrataHealthHandler(c *gin.Context) {
	ok, details := cache.CheckStrataKV()
	statusText := "ok"
	statusCode := http.StatusOK
	if !ok {
		statusText = "error"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":  statusText,
		"details": details,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}
