package lib

import (
	"log"

	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

type User struct {
	gorm.Model
	Name     string
	PassWord string
}

func Conn() *gorm.DB {
	var err error
	Db, err = gorm.Open("sqlite3", "./test.db")
	//Db, _ = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
	Db.AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
	}
	return Db
}
