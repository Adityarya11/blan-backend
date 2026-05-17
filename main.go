package main

import (
	"blan-backend/api"
	"blan-backend/runner"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/rs/cors"
)

func main() {

	runner.InitWorkerPool(3, 100)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	v1.POST("/compile", api.CompileHandler)
	v1.GET("/status/:id", api.StatusHandler)

	log.Println("Blan Backend is running on port 8080..")
	r.Run(":8080")
}
