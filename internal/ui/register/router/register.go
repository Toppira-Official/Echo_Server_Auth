package router

import (
	"auth/internal/ui/register/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRegisterRoutes(r *gin.Engine, c *controller.Register) {
	v1 := r.Group("/api/v1/auth")
	v1.POST("/register", c.Register)
}
