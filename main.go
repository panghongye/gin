package main

import (
	"gin/lib"
	"gin/middleware"
	"gin/model"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type (
	User model.User
	H    gin.H
)

var (
	db      *gorm.DB
	SendErr = lib.SendErr
)

func main() {
	{
		// 记录到文件。
		// f, _ := os.Create("gin.log")
		// gin.DefaultWriter = io.MultiWriter(f)
		// 如果需要同时将日志写入文件和控制台，请使用以下代码。
		// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	db = model.Conn()
	defer db.Close()
	ws := lib.GetWs()
	defer ws.Close()

	router := gin.Default()
	router.Use(gin.Recovery(), middleware.Cros)

	router.Static("/assets", "./assets")
	router.Static("/dist", "./dist")
	router.Any("/socket.io/*any", gin.WrapH(ws))

	{
		v1 := router.Group("/v1")
		v1.GET("/user/", GetUsers)
		v1.GET("/user/:id", GetUser)
		v1.POST("/user", AddUser)
		v1.PUT("/user/:id", UpdateUser)
		v1.DELETE("/user/:id", DeleteUser)
	}

	{
		router.GET("/alive", func(c *gin.Context) {
			c.JSON(200, H{
				"message": "server alive",
				"time":    time.Now(),
			})
		})
		v1 := router.Group("/api/v1")
		v1.POST("/register", func(c *gin.Context) {
			var t struct {
				Name     string
				Password string
			}
			if SendErr(c.Bind(&t), c) {
				return
			}
			log.Fatalln(t)
		})
		v1.POST("/login", func(c *gin.Context) {})
		v1.POST("/github_oauth", func(c *gin.Context) {})
	}

	router.Run(":3333")
}

func DeleteUser(c *gin.Context) {
	var t User
	if d := db.Where("id = ?", c.Params.ByName("id")).Delete(&t); d.Error == nil {
		c.JSON(200, H{"data": t, "d": d})
	} else {
		SendErr(d.Error, c)
	}
}

func UpdateUser(c *gin.Context) {
	var t User
	if d := db.Where("id = ?", c.Params.ByName("id")).First(&t); d.Error != nil {
		c.AbortWithStatus(404)
		log.Println("UpdateUser", d.Error)
	}
	if SendErr(c.BindJSON(&t), c) {
		return
	}
	if d := db.Save(&t); d.Error != nil {
		SendErr(d.Error, c)
	} else {
		c.JSON(200, H{
			"data": t,
			"d":    d,
		})
	}
}

func AddUser(c *gin.Context) {
	var t User
	if SendErr(c.BindJSON(&t), c) {
		return
	}
	if d := db.Create(&t); d.Error == nil {
		c.JSON(200, H{
			"data": t, "d": d,
		})
	} else {
		SendErr(d.Error, c)
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var t User
	if d := db.Where("id = ?", id).First(&t); d.Error != nil {
		c.AbortWithStatus(404)
		log.Println(d.Error)
	} else {
		c.JSON(200, H{
			"data": t, "d": d,
		})
	}
}

func GetUsers(c *gin.Context) {
	var t []User
	if d := db.Find(&t); d.Error != nil {
		c.AbortWithStatus(404)
		log.Println(d.Error)
	} else {
		c.JSON(200, H{"data": t, "d": d})
	}
}
