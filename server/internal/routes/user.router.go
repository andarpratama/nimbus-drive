package routes
import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)
func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", middleware.AuthMiddleware(), handlers.GetAllUsers)
	rg.GET("/users/:id", middleware.AuthMiddleware(), handlers.GetUserByID)
	rg.PUT("/users/profile", middleware.AuthMiddleware(), handlers.UpdateProfile)
	rg.PUT("/users/password", middleware.AuthMiddleware(), handlers.UpdatePassword)
}
