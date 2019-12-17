package service

import (
	"gin/model/table"
	"time"

	"github.com/jinzhu/gorm"
)

type GroupService struct{}

func (GroupService) FuzzyMatchGroups(name string) []table.GroupInfo {
	t := []table.GroupInfo{}
	sql := `SELECT * FROM group_info WHERE name LIKE ?;`
	db.Raw(sql, name).Scan(&t)
	return t
}

func (GroupService) JoinGroup(groupID string, ToUserIDs ...int) *gorm.DB {
	t := []table.GroupUserRelation{}
	for _, id := range ToUserIDs {
		t = append(t, table.GroupUserRelation{GroupID: groupID, UserID: id})
	}
	return db.Create(&t)
}

func (GroupService) IsInGroup(user_id int, group_id string) *gorm.DB {
	_sql := `SELECT * FROM group_user_relation WHERE user_id = ? AND group_id = ?;`
	return db.Raw(_sql, user_id, group_id)
}

func (GroupService) CreateGroup(name, group_notice, group_id string, from_user int) *gorm.DB {
	_sql :=
		`INSERT INTO group_info (id,name,group_notice,from_user,create_time) VALUES (?,?,?,?,?)`
	return db.Exec(_sql, group_id, name, group_notice, from_user, time.Now())
}

// 更新群信息
func (GroupService) UpdateGroupInfo(name, group_notice string, group_id string) *gorm.DB {
	var _sql = `UPDATE group_info SET name = ?, group_notice = ? WHERE group_id= ? limit 1 ; `
	return db.Exec(_sql, name, group_notice, group_id)
}

// 退群
func (GroupService) LeaveGroup(user_id, group_id string) *gorm.DB {
	var _sql = `DELETE FROM group_user_relation WHERE user_id = ? AND group_id = ? ;`
	return db.Exec(_sql, user_id, group_id)
}
