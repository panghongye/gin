package route

import (
	"gin/controller"
	"gin/lib"
	"gin/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	userCtrl controller.UserCtrl
)

func BuildRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Cros)
	{
		v1 := router.Group("/api/v1")
		v1.GET("/alive", func(c *gin.Context) {
			c.JSON(200, map[string]interface{}{
				"message": "server alive",
				"time":    time.Now(),
			})
		})
		v1.POST("/login", userCtrl.Login)
		v1.POST("/github_oauth")
		v1.POST("/register", userCtrl.Register)
	}
	{
		ws := lib.GetWs()
		router.Any("/socket.io/*any", gin.WrapH(ws))
	}

	return router
}
