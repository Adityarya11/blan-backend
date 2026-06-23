package api

import (
	"blan-backend/cache"
	"blan-backend/models"
	"blan-backend/runner"
	"blan-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompileHandler(g *gin.Context) {
	var req models.CompileRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	hashKey := utils.GenerateCacheKey(req.SourceCode)
	if cacheOutput, exists := cache.GetCachedOutput(hashKey); exists {
		g.JSON(http.StatusOK, models.CompileResponse{
			Output: cacheOutput,
			Cached: true,
		})
		return
	}

	jobID := runner.EnqueueJob(req.SourceCode, hashKey)

	g.JSON(http.StatusAccepted, models.CompileAcceptedResponse{
		ID:     jobID,
		Status: string(runner.JobQueued),
	})
}

func StatusHandler(g *gin.Context) {
	jobID := g.Param("id")

	record, ok := runner.GetJob(jobID)
	if !ok {
		g.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	g.JSON(http.StatusOK, models.JobStatusResponse{
		ID:     jobID,
		Status: string(record.Status),
		Output: record.Output,
		Error:  record.Error,
	})
}
