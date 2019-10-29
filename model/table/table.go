package table

import "github.com/jinzhu/gorm"

type UserInfo struct {
	gorm.Model
	Github_id uint   `json:"github_id"`
	Name      string `json:"name"`
	Password  string
	Avatar    string `json:"avatar"`
	Location  string `json:"location"`
	Socketid  string `json:"socketid"`
	Website   string `json:"website"`
	Github    string `json:"github"`
	Intro     string `json:"intro"`
	Company   string `json:"company"`
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
