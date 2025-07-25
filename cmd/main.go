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

	// ✅ AutoMigrate all models
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	api := r.Group("/api")
	routes.RegisterAuthRoutes(api)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello World"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
