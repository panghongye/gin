package service

import (
	"gin/model/table"
)

type UserService struct{}

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

// 新增用户
func (this UserService) InsertData(user *table.UserInfo) *table.UserInfo {
	db.Create(user)
	return user
}
