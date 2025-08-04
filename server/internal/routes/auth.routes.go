package routes
import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)
func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", handlers.Register)
	rg.POST("/login", handlers.Login)
	rg.POST("/logout", handlers.Logout)
	rg.GET("/user", middleware.AuthMiddleware(), handlers.GetCurrentUser)
}
