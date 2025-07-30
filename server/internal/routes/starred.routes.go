package routes

import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterStarredRoutes(starred *gin.RouterGroup) {
	protected := starred.Group("/starred")
	protected.Use(middleware.AuthMiddleware())

	// Star an item (file or folder)
	protected.POST("", handlers.StarItem)

	// Unstar an item (file or folder)
	protected.DELETE("", handlers.UnstarItem)

	// Toggle star status for an item (file or folder)
	protected.PATCH("", handlers.ToggleStarItem)

	// List all starred items
	protected.GET("", handlers.ListStarredItems)
} 