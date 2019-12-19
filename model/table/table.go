package table

import "time"

type UserInfo struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Name     string `gorm:"unique;not null" json:"name"`
	Password string `json:"-"`
	Intro    string `json:"intro"`
}

type GroupInfo struct {
	ID         string    `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	Intro      string    `json:"intro"`
	FromUser   int       `json:"fromUser"`
	CreateTime time.Time `json:"createTime"`
}

type GroupMsg struct {
	ID          int       `gorm:"primary_key" json:"id"`
	FromUser    int       `json:"fromUser"`
	GroupID     string    `json:"groupId"`
	Msg         string    `json:"message"`
	Time        time.Time `json:"time"`
	Attachments string    `gorm:"type:json" json:"attachments"`
}

type GroupUserRelation struct {
	ID      int    `gorm:"primary_key" json:"id"`
	GroupID string `gorm:"UNIQUE_INDEX:GroupUser" json:"groupId"`
	UserID  int    `gorm:"UNIQUE_INDEX:GroupUser" json:"userID"`
}
