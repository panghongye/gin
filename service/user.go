package service

import (
	"gin/model/table"

	"github.com/jinzhu/gorm"
)

type (
	UserInfo table.UserInfo
)

func (_ UserInfo) Login(name string) *gorm.DB {
	return DB.Where("name LIKE ?", name).Find(&UserInfo{})
}

// 模糊匹配用户
func (_ UserInfo) FuzzyMatchUsers(name string) *gorm.DB {
	return DB.Where("name LIKE ?", name).Find(&UserInfo{})
}

// 通过用户名查找非github用户信息 user_info
func (_ UserInfo) FindDataByName(name string) *gorm.DB {
	return DB.Where("name = ?", name).First(&UserInfo{})
}

// 注册用户
func (_ UserInfo) InsertData(user *UserInfo) *gorm.DB {
	return DB.Create(user)
}
