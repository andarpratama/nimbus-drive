package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	RegisterAuthRoutes(api)

	// You can register other route groups here, e.g.
	// RegisterFileRoutes(api)
	// RegisterFolderRoutes(api)
}
