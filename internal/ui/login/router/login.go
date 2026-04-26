package router

import (
	"auth/internal/ui/login/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthLoginRoutes(r *gin.Engine, c *controller.Login) {
	v1 := r.Group("/api/v1/auth")
	v1.POST("/login", c.Login)
}
