package database

import (
	"blan-backend/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dbconnect string) {
	var err error

	DB, err = gorm.Open(mysql.Open(dbconnect), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}

	log.Println("Connected to the SQL")

	err = DB.AutoMigrate(&models.User{}, &models.Snippet{}, &models.Job{})
	if err != nil {
		log.Fatalf("Failed to db migrations %v", err)
	}

	log.Println("DB migration successful.")
}
