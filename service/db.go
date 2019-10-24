package service

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	DB, err := gorm.Open("sqlite3", "./test.db")
	//DB, err := gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	// DB.AutoMigrate(&User{}, &UserInfo{}, &Group_info{}, &User_user_relation{}, &Group_msg{}, &Private_msg{}, &Group_user_relation{})

	DB.LogMode(true)
	DB.SingularTable(true)
	// return DB
}
