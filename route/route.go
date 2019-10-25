package route

import (
	"gin/controller"
	"gin/middleware"

	"github.com/gin-gonic/gin"
)

var (
	userCtrl controller.UserCtrl
)

func BuildRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), middleware.Cros)
	v1 := router.Group("api/v1")
	{
		v1.POST("/login", userCtrl.Login)
		v1.POST("/github_oauth")
		v1.POST("/register", userCtrl.Register)
	}
	return router
}
