package routes
import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)
func RegisterStarredRoutes(starred *gin.RouterGroup) {
	protected := starred.Group("/starred")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("", handlers.StarItem)
	protected.DELETE("", handlers.UnstarItem)
	protected.PATCH("", handlers.ToggleStarItem)
	protected.GET("", handlers.ListStarredItems)
} 