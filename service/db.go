package service

import (
	"gin/model/table"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/spf13/viper"
)

var (
	db *gorm.DB
)

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "./test.db")
	// db, err = gorm.Open("mysql", viper.GetString("db.connStr"))
	if err != nil {
		panic(err.Error())
	}
	db.LogMode(true)
	db.SingularTable(true) // 关闭复数表名，如果设置为true，`User`表的表名就会是`user`，而不是`users`
	db.AutoMigrate(
		new(table.UserInfo),
		new(table.GroupInfo),
		new(table.GroupMsg),
		new(table.GroupUserRelation),
	)
}
