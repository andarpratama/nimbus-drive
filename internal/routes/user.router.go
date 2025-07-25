package routes

import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/andarpratama/nimbus-drive/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "superadmin"),
	)

	users.GET("", handlers.GetAllUsers)
	users.GET("/:id", handlers.GetUserByID)
}
