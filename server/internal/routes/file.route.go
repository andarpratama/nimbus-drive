package routes

import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(files *gin.RouterGroup) {
	protected := files.Group("/files")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/upload", handlers.UploadFile)
	protected.GET("/:id/download", handlers.DownloadFile)
	protected.DELETE("/:id", handlers.DeleteFile)
	protected.GET("", handlers.ListFiles)
	protected.GET("/trash", handlers.GetTrashedFiles)
	protected.PUT("/:id/move", handlers.MoveFile)
	protected.POST("/:id/restore", handlers.RestoreFile)
	protected.DELETE("/:id/permanent", handlers.PermanentlyDeleteFile)

	// Public route for serving images (no authentication required)
	files.Group("/files").GET("/:id/image", handlers.ServeImage)
}
