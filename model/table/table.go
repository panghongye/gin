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
	IsFriend   int       `gorm:"DEFAULT:0" json:"isFriend"`
}

type GroupMsg struct {
	ID          int       `gorm:"primary_key" json:"id"`
	FromUser    int       `json:"fromUser"`
	GroupID     string    `json:"groupID"`
	Msg         string    `gorm:"Type:text" json:"msg"`
	Time        time.Time `json:"time"`
	Attachments string    `gorm:"type:json" json:"attachments"`
}

type GroupUserRelation struct {
	ID      int    `gorm:"primary_key" json:"id"`
	GroupID string `gorm:"UNIQUE_INDEX:GroupUser" json:"groupID"`
	UserID  int    `gorm:"UNIQUE_INDEX:GroupUser" json:"userID"`
}
