package table

type UserInfo struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Github_id int    `json:"github_id"`
	Name      string `json:"name"`
	Password  string `json:"-"`
	Avatar    string `json:"avatar"`
	Location  string `json:"location"`
	Socketid  string `json:"socketid"`
	Website   string `json:"website"`
	Github    string `json:"github"`
	Intro     string `json:"intro"`
	Company   string `json:"company"`
}

type Group_info struct {
	ID           int    `gorm:"primary_key" json:"id"`
	To_group_id  string `json:"to_group_id"`
	Name         string `json:"name"`
	Group_notice string `json:"group_notice"`
	Creator_id   int    `json:"creator_id"`
	Create_time  int    `json:"create_time"`
}

type User_user_relation struct {
	ID        int    `gorm:"primary_key" json:"id"`
	User_id   int    `json:"user_id"`
	From_user int    `json:"from_user"`
	Remark    string `json:"remark"`
	Shield    uint64 `json:"shield"`
	Time      int    `json:"time"`
}
type Group_msg struct {
	ID          int    `gorm:"primary_key" json:"id"`
	From_user   int    `json:"from_user"`
	To_group_id string `json:"to_group_id"`
	Message     string `json:"message"`
	Time        int    `json:"time"`
	Attachments string `json:"attachments"`
}

type Private_msg struct {
	ID          int    `gorm:"primary_key" json:"id"`
	From_user   int    `json:"from_user"`
	To_user     int    `json:"to_user"`
	Message     string `json:"message"`
	Time        int    `json:"time"`
	Attachments string `json:"attachments"`
}

type Group_user_relation struct {
	ID          int    `gorm:"primary_key" json:"id"`
	To_group_id string `json:"to_group_id"`
	User_id     string `json:"user_id"`
}
