package service

import (
	"time"

	"github.com/jinzhu/gorm"
)

// "gin/model/table"

// "github.com/jinzhu/gorm"

type GroupChatService struct{}

func (GroupChatService) GetGroupMsg(to_group_id string, start, count int) *gorm.DB {
	sql := `SELECT * FROM (SELECT g.message,g.attachments,g.time,g.from_user,g.to_group_id, i.avatar ,i.name, i.github_id FROM group_msg  As g inner join user_info AS i ON g.from_user = i.id  WHERE to_group_id = ? order by time desc limit ?,?) as n order by n.time; `
	return db.Raw(sql, to_group_id, start, count)
}

func (GroupChatService) GetGroupMember(to_group_id string, start, count int) *gorm.DB {
	sql := `SELECT g.user_id, u.socketid, u.name, u.avatar, u.github_id, u.github, u.intro, u.company, u.location, u.website FROM group_user_relation AS g inner join user_info AS u ON g.user_id = u.id WHERE to_group_id = ?`
	return db.Raw(sql, to_group_id)
}

func (GroupChatService) GetGroupInfo(to_group_id string, name string) *gorm.DB {
	return db.Raw(`SELECT to_group_id, name, group_notice, creator_id, create_time FROM group_info  WHERE to_group_id = ? OR name = ? ;`, to_group_id, name)
}

func (GroupChatService) SaveGroupMsg(from_user int, to_group_id, message, attachments string) *gorm.DB {
	sql := `INSERT group_msg(from_user,to_group_id,message ,time, attachments) VALUES(?,?,?,?,?); `
	return db.Exec(sql, from_user, to_group_id, message, int(time.Now().Unix()), attachments)
}

func (GroupChatService) AddGroupUserRelation(user_id, groupId int) *gorm.DB {
	const _sql = `INSERT INTO  group_user_relation(to_group_id,user_id) VALUES(?,?); `
	return db.Raw(_sql, groupId, user_id)
}

func (GroupChatService) GetUnreadCount(sortTime int, to_group_id string) *gorm.DB {
	return db.Raw(`SELECT count(time) as unread FROM group_msg as p where p.time > ? and p.to_group_id = ?;`, sortTime, to_group_id)
}
