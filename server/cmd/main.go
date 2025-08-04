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
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}
	if err := database.ConnectMySQL(); err != nil {
		log.Fatal("DB error:", err)
	}
	database.ConnectRedis()
	resetDB := os.Getenv("RESET_DB")
	log.Printf("🔍 RESET_DB environment variable: '%s'", resetDB)
	if resetDB == "true" {
		log.Println("🔄 Resetting database...")
		if err := database.DB.AutoMigrate(
			&models.User{},
			&models.File{},
			&models.Folder{},
			&models.SharedFile{},
			&models.Session{},
			&models.Starred{},
		); err != nil {
			log.Fatal("Migration error:", err)
		}
		log.Println("✅ Database reset successful")
	} else {
		log.Println("📊 Database reset disabled (RESET_DB != true)")
		if err := database.DB.AutoMigrate(
			&models.User{},
			&models.File{},
			&models.Folder{},
			&models.SharedFile{},
			&models.Session{},
			&models.Starred{},
		); err != nil {
			log.Printf("⚠️ Migration warning (non-critical): %v", err)
		}
	}
	r := gin.Default()
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, 
	}
	r.Use(cors.New(corsConfig))
	routes.RegisterRoutes(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello World"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server starting on port %s", port)
	r.Run(":" + port)
}
