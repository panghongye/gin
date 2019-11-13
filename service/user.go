package service

import (
	"gin/model/table"

	"github.com/jinzhu/gorm"
)

type UserService struct{}

// 新增用户
func (this UserService) InsertData(user *table.UserInfo) *table.UserInfo {
	db.Create(user)
	return user
}

// 模糊匹配用户
func (this UserService) FuzzyMatchUsers(name string) *table.UserInfo {
	t := new(table.UserInfo)
	db.Where("name LIKE ?", name).Find(t)
	return t
}

// 通过用户名查用户
func (this UserService) FindDataByName(name string) *table.UserInfo {
	t := new(table.UserInfo)
	db.Where("name = ?", name).First(t)
	return t
}

func (this UserService) GetByID(userID int) *table.UserInfo {
	t := new(table.UserInfo)
	db.Where("id = ?", userID).First(t)
	return t
}
func (this UserService) SaveUserSocketId(userID int, socketId string) *gorm.DB {
	return db.Model(new(table.UserInfo)).Where("id = ?", userID).Update("socketid", socketId)
}
