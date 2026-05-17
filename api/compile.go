package api

import (
	"blan-backend/models"
	"blan-backend/runner"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompileHandler(g *gin.Context) {
	var req models.CompileRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request Body"})
		return
	}

	jobID := runner.EnqueueJob(req.SourceCode)

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
