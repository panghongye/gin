package route

import (
	// "gin/controller"
	"gin/middleware"

	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), middleware.Cros)
	v1 := router.Group("api/v1")
	{
		v1.POST("/login",)
	}
	return router
}
