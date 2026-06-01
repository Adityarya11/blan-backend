package main

import (
	"blan-backend/api"
	"blan-backend/cache"
	"blan-backend/database"
	"blan-backend/middleware"
	"blan-backend/runner"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		// log.Fatalf("Error with the env: %v", err) // wont work with the docker
		log.Println("No .env file found, relying on the system environments.")
	}

	dbconnect := os.Getenv("DATABASE_URL")
	database.Connect(dbconnect)

	cache.InitStrataKV("./strata_cache_data")
	defer cache.CloseStrata()

	runner.InitWorkerPool(3, 100)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	v1.POST("/compile", api.CompileHandler)
	v1.GET("/status/:id", api.StatusHandler)
	v1.GET("/health/strata", api.StrataHealthHandler)

	v1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "alive",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	//auth routes
	v1.POST("/signup", api.SignupHandler)
	v1.POST("/login", api.LoginHandler)

	secured := v1.Group("/snippets")
	secured.Use(middleware.AuthMiddleware())
	{
		secured.POST("/", api.CreateSnippetHandler)
		secured.GET("/", api.GetSnippetHandler)
	}

	log.Println("Blan Backend is running on port 8080..")
	r.Run(":8080")
}
