package model

import (
	"log"

	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type Room struct {
	gorm.Model
	Name string
}

type RoomUser struct {
	gorm.Model
	RoomId uint
	UserId uint
}

func Conn() *gorm.DB {
	var err error
	Db, err = gorm.Open("sqlite3", "./test.db")
	//Db, err = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
	Db.AutoMigrate(&User{}, &Room{}, &RoomUser{})
	if err != nil {
		log.Println(err)
	}
	return Db
}
