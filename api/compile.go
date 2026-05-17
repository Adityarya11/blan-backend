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

	output, err := runner.RunSource(req.SourceCode)

	if err != nil {
		g.JSON(http.StatusOK, models.CompileResponse{
			Output: "",
			Error:  err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, models.CompileResponse{
		Output: output,
		Error:  "",
	})
}
