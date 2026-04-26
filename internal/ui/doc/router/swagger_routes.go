package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterSwaggerRoutes(r *gin.Engine) {
	g := r.Group("/docs/swagger")
	g.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
