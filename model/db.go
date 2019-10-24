package model

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type UserInfo struct {
	gorm.Model
	Github_id uint
	Name      string
	Password  string
	Avatar    string
	Location  string
	Socketid  string
	Website   string
	Github    string
	Intro     string
	Company   string
}

type Group_info struct {
	gorm.Model
	To_group_id  string
	Name         string
	Group_notice string
	Creator_id   uint
}

type User_user_relation struct {
	gorm.Model
	User_id   uint
	From_user uint
	Remark    string
	Shield    uint64
	Time      uint
}
type Group_msg struct {
	gorm.Model
	From_user   uint
	To_group_id string
	Message     string
	Time        uint
	Attachments string
}

type Private_msg struct {
	gorm.Model
	From_user   uint
	To_user     uint
	Message     string
	Time        uint
	Attachments string
}
type Group_user_relation struct {
	gorm.Model
	To_group_id string
	User_id     string
}

func Conn() *gorm.DB {
	var err error
	DB, err = gorm.Open("sqlite3", "./test.db")
	//DB, err = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
	DB.AutoMigrate(&User{}, &UserInfo{}, &Group_info{}, &User_user_relation{}, &Group_msg{}, &Private_msg{}, &Group_user_relation{})
	if err != nil {
		log.Println(err)
	}
	return DB
}
