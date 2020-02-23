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
	db.Raw(`SELECT * FROM group_info WHERE id in (SELECT group_id  FROM group_user_relation WHERE user_id=?) ORDER BY create_time`, useID).Scan(&t)
	return t
}

type GroupMsg struct {
	table.GroupMsg
	UserName string `json:"userName"`
}

func (GroupService) FindGroupMsgByGroupID(groupID string, page, pageSize int) []GroupMsg {
	pageParam(&page, &pageSize)
	t := []GroupMsg{} // LIMIT ?,?  (page-1)*pageSize, pageSize
	db.Raw(`SELECT G.*,U.name as user_name FROM group_msg AS G LEFT JOIN user_info AS U ON G.from_user=U.id WHERE group_id=?`, groupID).Scan(&t)
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

func (GroupService) CreateGroup(name, Intro, group_id string, from_user int, is_friend int) *gorm.DB {
	_sql :=
		`INSERT INTO group_info (id,name,Intro,from_user,create_time,is_friend) VALUES (?,?,?,?,?,?)`
	return db.Exec(_sql, group_id, name, Intro, from_user, time.Now(), is_friend)
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

func pageParam(page *int, pageSize *int) {
	if *page < 1 {
		*page = 1
	}

	if *pageSize < 1 {
		*pageSize = 1
	}
}
