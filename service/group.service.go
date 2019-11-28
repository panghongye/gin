package service

import (
	"gin/model/table"
	"time"

	"github.com/jinzhu/gorm"
)

type GroupService struct{}

func (GroupService) FuzzyMatchGroups(name string) []table.Group_info {
	t := []table.Group_info{}
	sql := `SELECT * FROM group_info WHERE name LIKE ?;`
	db.Raw(sql, name).Scan(&t)
	return t
}

func (GroupService) JoinGroup(user_id int, to_group_id string) *gorm.DB {
	_sql := `INSERT INTO group_user_relation(user_id,to_group_id) VALUES(?,?);`
	return db.Exec(_sql, user_id, to_group_id)
}

func (GroupService) IsInGroup(user_id int, to_group_id string) *gorm.DB {
	_sql := `SELECT * FROM group_user_relation WHERE user_id = ? AND to_group_id = ?;`
	return db.Exec(_sql, user_id, to_group_id)
}

func (GroupService) CreateGroup(name, group_notice, to_group_id string, creator_id int) *gorm.DB {
	_sql :=
		`INSERT group_info (to_group_id,name,group_notice,creator_id,create_time) VALUES (?,?,?,?,?)`
	return db.Exec(_sql, to_group_id, name, group_notice, creator_id, time.Now())
}

// 更新群信息
func (GroupService) UpdateGroupInfo(name, group_notice string, to_group_id string) *gorm.DB {
	var _sql = `UPDATE group_info SET name = ?, group_notice = ? WHERE to_group_id= ? limit 1 ; `
	return db.Exec(_sql, name, group_notice, to_group_id)
}

// 退群
func (GroupService) LeaveGroup(user_id, to_group_id string) *gorm.DB {
	var _sql = `DELETE FROM group_user_relation WHERE user_id = ? AND to_group_id = ? ;`
	return db.Exec(_sql, user_id, to_group_id)
}
