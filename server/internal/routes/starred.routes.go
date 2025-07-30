package routes

import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterStarredRoutes(starred *gin.RouterGroup) {
	protected := starred.Group("/starred")
	protected.Use(middleware.AuthMiddleware())

	// Star/unstar files
	protected.POST("/files/:id/star", handlers.StarFile)
	protected.DELETE("/files/:id/star", handlers.UnstarFile)

	// Star/unstar folders
	protected.POST("/folders/:id/star", handlers.StarFolder)
	protected.DELETE("/folders/:id/star", handlers.UnstarFolder)

	// List all starred items
	protected.GET("", handlers.ListStarredItems)

	// Check if an item is starred
	protected.GET("/:id/status", handlers.CheckStarredStatus)
} 