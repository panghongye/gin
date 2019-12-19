package service

import (
	"gin/model/table"
	"time"

	"github.com/jinzhu/gorm"
)

type GroupService struct{}

func (GroupService) FuzzyFindGroupsByName(name string) *[]table.GroupInfo {
	t := []table.GroupInfo{}
	sql := `SELECT * FROM group_info WHERE name LIKE ?;`
	db.Raw(sql, name).Scan(&t)
	return &t
}

func (GroupService) FindGroupByID(id string) table.GroupInfo {
	t := table.GroupInfo{ID: id}
	db.First(&t)
	return t
}

//FindGroupsByUserID 获取用户所在群组
func (GroupService) FindGroupsByUserID(useID int) []table.GroupInfo {
	t := []table.GroupInfo{}
	db.Raw(`SELECT * FROM group_info WHERE id in (SELECT group_id  FROM group_user_relation WHERE user_id=?)`, useID).Scan(&t)
	return t
}

func (GroupService) FindGroupMsgByGroupID(groupID string, page, pageSize int) []table.GroupMsg {
	t := []table.GroupMsg{}
	db.Raw(`SELECT * FROM group_msg WHERE group_id='?' ORDER BY id DESC LIMIT ?,?`, groupID, (page-1)*pageSize, pageSize).Scan(&t)
	return t
}

func (GroupService) JoinGroup(groupID string, UserIDs ...int) {
	for _, id := range UserIDs {
		db.Create(&table.GroupUserRelation{GroupID: groupID, UserID: id})
	}
}

func (GroupService) IsInGroup(user_id int, group_id string) *gorm.DB {
	_sql := `SELECT * FROM group_user_relation WHERE user_id = ? AND group_id = ?;`
	return db.Raw(_sql, user_id, group_id)
}

func (GroupService) CreateGroup(name, Intro, group_id string, from_user int) *gorm.DB {
	_sql :=
		`INSERT INTO group_info (id,name,Intro,from_user,create_time) VALUES (?,?,?,?,?)`
	return db.Exec(_sql, group_id, name, Intro, from_user, time.Now())
}

// 更新群信息
func (GroupService) UpdateGroupInfo(name, Intro string, group_id string) *gorm.DB {
	var _sql = `UPDATE group_info SET name = ?, Intro = ? WHERE group_id= ? limit 1 ; `
	return db.Exec(_sql, name, Intro, group_id)
}

// 退群
func (GroupService) LeaveGroup(user_id, group_id string) *gorm.DB {
	var _sql = `DELETE FROM group_user_relation WHERE user_id = ? AND group_id = ? ;`
	return db.Exec(_sql, user_id, group_id)
}
