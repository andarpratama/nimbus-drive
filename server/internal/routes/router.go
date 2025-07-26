package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	RegisterAuthRoutes(api)
	RegisterFileRoutes(api)
	RegisterUserRoutes(api)
}
