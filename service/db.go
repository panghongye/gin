package service

import (
	"gin/model/table"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "./test.db")
	//DB, err = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	db.LogMode(true)
	db.SingularTable(true) // 关闭复数表名，如果设置为true，`User`表的表名就会是`user`，而不是`users`
	db.AutoMigrate(
		new(table.UserInfo),
		new(table.User_user_relation),
		new(table.Group_msg),
		new(table.Group_info),
		new(table.Group_user_relation),
		new(table.Private_msg),
	)
}
