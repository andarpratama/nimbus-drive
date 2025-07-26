package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/andarpratama/nimbus-drive/internal/database"
	"github.com/andarpratama/nimbus-drive/internal/models"
	"github.com/andarpratama/nimbus-drive/internal/routes"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

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

	// Create Gin router
	r := gin.Default()

	// Configure CORS - Apply before any routes
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://127.0.0.1:5173", "http://127.0.0.1:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	r.Use(cors.New(corsConfig))

	// Register API routes
	routes.RegisterRoutes(r)

	// Optional: health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello World"})
	})

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server starting on port %s", port)
	r.Run(":" + port)
}
