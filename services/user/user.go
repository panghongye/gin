package user

import (
	"gin/lib"
	"gin/model"

	"github.com/jinzhu/gorm"
)

var (
	DB      = model.DB
	SendErr = lib.SendErr
)

type (
	UserInfo model.UserInfo
)

// 模糊匹配用户
func FuzzyMatchUsers(name string) *gorm.DB {
	return DB.Where("name LIKE ?", name).Find(&UserInfo{})
}

// 通过用户名查找非github用户信息 user_info
func FindDataByName(name string) *gorm.DB {
	return DB.Where("name = ?", name).First(&UserInfo{})
}

// 注册用户
func InsertData(user *UserInfo) *gorm.DB {
	return DB.Create(user)
}
