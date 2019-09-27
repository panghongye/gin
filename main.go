package main

import (
	"gin/lib"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type (
	User lib.User
	H    gin.H
)

var (
	db      *gorm.DB
	SendErr = lib.SendErr
)

func main() {
	db = lib.Conn()
	defer db.Close()
	r := gin.Default()
	r.GET("/people/", GetPeople)
	r.GET("/people/:id", GetUser)
	r.POST("/people", AddUser)
	r.PUT("/people/:id", UpdateUser)
	r.DELETE("/people/:id", DeleteUser)
	r.Run(":8080")
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
