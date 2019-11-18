package service

import "github.com/jinzhu/gorm"

type GroupService struct{}

func (GroupService) FuzzyMatchGroups(name string) *gorm.DB {
	sql := `SELECT * FROM group_info WHERE name LIKE ?;`
	return db.Raw(sql, name)
}

func (GroupService) JoinGroup(user_id uint, to_group_id string) *gorm.DB {
	_sql := `INSERT INTO group_user_relation(user_id,to_group_id) VALUES(?,?);`
	return db.Raw(_sql, user_id, to_group_id)
}

// 查看某个用户是否在某个群中
func (GroupService) IsInGroup(user_id, to_group_id string) *gorm.DB {
	_sql := `SELECT * FROM group_user_relation WHERE user_id = ? AND to_group_id = ?;`
	return db.Raw(_sql, user_id, to_group_id)
}

// 建群
func (GroupService) CreateGroup(name, group_notice, to_group_id string, creator_id uint, create_time int64) *gorm.DB {
	_sql :=
		`INSERT INTO group_info (to_group_id,name,group_notice,creator_id,create_time) VALUES (?,?,?,?,?)`
	return db.Raw(_sql, to_group_id, name, group_notice, creator_id, create_time)
}

// 更新群信息
func (GroupService) UpdateGroupInfo(name, group_notice string, to_group_id string) *gorm.DB {
	var _sql = `UPDATE group_info SET name = ?, group_notice = ? WHERE to_group_id= ? limit 1 ; `
	return db.Raw(_sql, name, group_notice, to_group_id)
}

// 退群
func (GroupService) LeaveGroup(user_id, to_group_id string) *gorm.DB {
	var _sql = `DELETE FROM group_user_relation WHERE user_id = ? AND to_group_id = ? ;`
	return db.Raw(_sql, user_id, to_group_id)
}
