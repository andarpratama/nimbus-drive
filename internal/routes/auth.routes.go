package routes

import (
	"github.com/andarpratama/nimbus-drive/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", handlers.Register)
	rg.POST("/login", handlers.Login)
}
