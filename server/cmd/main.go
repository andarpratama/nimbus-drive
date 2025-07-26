package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models" // ✅ make sure this matches your path
	"github.com/andarpratama/nimbus-drive/internal/routes"
)

func main() {
	godotenv.Load()

	if err := database.ConnectMySQL(); err != nil {
		log.Fatal("DB error:", err)
	}
	database.ConnectRedis()

	if err := database.DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Folder{},
		&models.SharedFile{},
		&models.Session{},
	); err != nil {
		log.Fatal("Migration error:", err)
	}
	log.Println("✅ AutoMigrate successful")

	// ❌ FIXED: only one instance of gin.Default
	r := gin.Default()

	// Register API routes
	routes.RegisterRoutes(r)

	// Optional: health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello World"})
	})

	err := godotenv.Load("../../.env") // relative path from server/main.go
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
