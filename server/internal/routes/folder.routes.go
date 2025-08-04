package routes
import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)
func RegisterFolderRoutes(folders *gin.RouterGroup) {
	protected := folders.Group("/folders")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("", handlers.CreateFolder)
	protected.GET("", handlers.GetFolders)
	protected.GET("/tree", handlers.GetFolderTree)
	protected.GET("/:id", handlers.GetFolderByID)
	protected.GET("/:id/contents", handlers.GetFolderContents)
	protected.PUT("/:id", handlers.UpdateFolder)
	protected.PATCH("/:id/rename", handlers.RenameFolder)
	protected.PATCH("/:id/move", handlers.MoveFolder)
	protected.DELETE("/:id", handlers.DeleteFolder)
} 