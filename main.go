package main

import (
	"gin/lib"
	"gin/model"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	// 记录到文件。
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	db = model.Conn()
	defer db.Close()

	ws := lib.GetWs()
	defer ws.Close()

	router := gin.Default()
	router.Static("/assets", "./assets")
	router.GET("/people/", GetPeople)
	router.GET("/people/:id", GetUser)
	router.POST("/people", AddUser)
	router.PUT("/people/:id", UpdateUser)
	router.DELETE("/people/:id", DeleteUser)
	router.Any("/socket.io/*any", gin.WrapH(ws))
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

func GetPeople(c *gin.Context) {
	var t []User
	if d := db.Find(&t); d.Error != nil {
		c.AbortWithStatus(404)
		log.Println(d.Error)
	} else {
		c.JSON(200, H{"data": t, "d": d})
	}
}
