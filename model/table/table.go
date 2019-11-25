package table

import "github.com/jinzhu/gorm"

type UserInfo struct {
	gorm.Model
	Github_id int   `json:"github_id"`
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
	Creator_id   int
}

type User_user_relation struct {
	gorm.Model
	User_id   int
	From_user int
	Remark    string
	Shield    uint64
	Time      int
}
type Group_msg struct {
	gorm.Model
	From_user   int
	To_group_id string
	Message     string
	Time        int
	Attachments string
}

type Private_msg struct {
	gorm.Model
	From_user   int
	To_user     int
	Message     string
	Time        int
	Attachments string
}

type Group_user_relation struct {
	gorm.Model
	To_group_id string
	User_id     string
}
