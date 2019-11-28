package service

import (
	"gin/model/table"
	"github.com/jinzhu/gorm"
)

type UserService struct{}

func (this UserService) FuzzyMatchUsers(name string) []table.UserInfo {
	t := []table.UserInfo{}
	db.Table("user_info").Where("name LIKE ?", name).Scan(&t)
	return t
}

func (this UserService) InsertData(user *table.UserInfo) *table.UserInfo {
	db.Create(user)
	return user
}

func (UserService) InsertGithubData(name, github_id, avatar, location, website, github, intro, company string) *gorm.DB {
	var _sql = `insert into user_info(name, github_id, avatar, location, website, github, intro, company) values(?,?,?,?,?,?,?,?);`
	return db.Raw(_sql, name, github_id, avatar, location, website, github, intro, company)
}

func (UserService) FindGithubUser(githubId string) *gorm.DB {
	var _sql = `SELECT * FROM user_info WHERE github_id = ? ;`
	return db.Raw(_sql, githubId)
}

func (UserService) UpdateGithubUser(name, avatar, location, website, github, intro, github_id, company string) *gorm.DB {
	var _sql = ` UPDATE  user_info SET name = ?,avatar = ?,location = ?,website = ?,github = ?,intro= ?, company = ? WHERE github_id = ? ; `
	return db.Raw(_sql, name, avatar, location, website, github, intro, company, github_id)
}

func (this UserService) FindDataByName(name string) *table.UserInfo {
	t := new(table.UserInfo)
	db.Where("name = ?", name).First(t)
	return t
}

func (UserService) GetUserInfo(user_id int) *gorm.DB {
	_sql := `SELECT id AS user_id, name, avatar, location, website, github, github_id, intro, company FROM user_info WHERE user_info.id =? `
	return db.Raw(_sql, user_id)
}

// 通过要查看的用户id 查询是否是本机用户的好友  如果是 返回user_id 和 remark 备注
func (UserService) IsFriend(user_id, from_user int) *gorm.DB {
	var _sql = `SELECT  * FROM user_user_relation  AS u WHERE  u.user_id = ? AND u.from_user = ? `
	return db.Raw(_sql, user_id, from_user)
}

// 两边都互加为好友
func (UserService) AddFriendEachOther(user_id, from_user, time int64) *gorm.DB {
	var _sql = `INSERT INTO user_user_relation(user_id,from_user,time) VALUES (?,?,?), (?,?,?)`
	return db.Raw(_sql, user_id, from_user, time, from_user, user_id, time)
}

// 删除联系人
func (UserService) DeleteContact(user_id, from_user int) *gorm.DB {
	var _sql = `DELETE FROM  user_user_relation WHERE (user_id = ? AND from_user = ?) or (user_id = ? AND from_user = ?)`
	return db.Raw(_sql, user_id, from_user, from_user, user_id)
}

// 通过user_id查找首页群列表
func (UserService) GetGroupList(user_id int) *gorm.DB {
	var _sql = `SELECT r.to_group_id ,i.name , i.create_time,
      (SELECT g.message  FROM group_msg AS g  WHERE g.to_group_id = r.to_group_id  ORDER BY TIME DESC   LIMIT 1 )  AS message ,
      (SELECT g.time  FROM group_msg AS g  WHERE g.to_group_id = r.to_group_id  ORDER BY TIME DESC   LIMIT 1 )  AS time,
      (SELECT g.attachments FROM group_msg AS g  WHERE g.to_group_id = r.to_group_id  ORDER BY TIME DESC   LIMIT 1 )  AS attachments
      FROM  group_user_relation AS r inner join group_info AS i on r.to_group_id = i.to_group_id WHERE r.user_id = ? ;`
	return db.Raw(_sql, user_id)
}

// 通过user_id查找首页私聊列表
func (UserService) GetPrivateList(user_id int) *gorm.DB {
	var _sql = ` SELECT r.from_user as user_id, i.name, i.avatar, i.github_id, r.time as be_friend_time,
      (SELECT p.message FROM private_msg AS p WHERE (p.to_user = r.from_user and p.from_user = r.user_id) or (p.from_user = r.from_user and p.to_user = r.user_id) ORDER BY p.time DESC   LIMIT 1 )  AS message ,
      (SELECT p.time FROM private_msg AS p WHERE (p.to_user = r.from_user and p.from_user = r.user_id) or (p.from_user = r.from_user and p.to_user = r.user_id) ORDER BY p.time DESC   LIMIT 1 )  AS time,
      (SELECT p.attachments FROM private_msg AS p WHERE (p.to_user = r.from_user and p.from_user = r.user_id) or (p.from_user = r.from_user and p.to_user = r.user_id) ORDER BY p.time DESC   LIMIT 1 )  AS attachments
      FROM  user_user_relation AS r inner join user_info AS i on r.from_user  = i.id WHERE r.user_id = ? ;`
	return db.Raw(_sql, user_id)
}

func (UserService) SaveUserSocketId(user_id int, socketId string) *gorm.DB {
	var _sql = ` UPDATE  user_info SET socketid = ? WHERE id= ? limit 1 ; `
	return db.Raw(_sql, socketId, user_id)
}

func (UserService) GetUserSocketId(toUserId int) *gorm.DB {
	var _sql = ` SELECT socketid FROM user_info WHERE id=? limit 1 ;`
	return db.Raw(_sql, toUserId)
}
