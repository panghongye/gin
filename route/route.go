package route

import (
	"gin/controller"
	"gin/lib/jwt"
	"gin/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	userCtrl controller.UserCtrl
)

func BuildRouter() *gin.Engine {
	jwt.New("0", time.Hour)
	router := gin.New()
	router.Static("/test", "./assets")
	router.Use(middleware.Cros)
	{
		v1 := router.Group("/api/v1")
		v1.POST("/login", userCtrl.Login)
		v1.POST("/register", userCtrl.Register)
		v1.POST("/github_oauth")
		router.Any("/socket.io/*any", middleware.Auth, gin.WrapH(getWs()))
		// router.Any("/socket.io/*any", gin.WrapH(getWs()))
	}
	return router
}
