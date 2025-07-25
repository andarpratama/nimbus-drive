package routes

import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(files *gin.RouterGroup) {
	protected := files.Group("/files")
	protected.Use(middleware.AuthMiddleware())

	files.POST("/upload", handlers.UploadFile)
	files.GET("/:id/download", handlers.DownloadFile)
	files.GET("", handlers.ListFiles)
}
