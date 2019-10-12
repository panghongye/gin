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

type ChatRoom struct {
	gorm.Model
	UserId uint
}

func Conn() *gorm.DB {
	var err error
	Db, err = gorm.Open("sqlite3", "./test.db")
	//Db, _ = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
	Db.AutoMigrate(&User{}, &ChatRoom{})
	if err != nil {
		log.Println(err)
	}
	return Db
}
