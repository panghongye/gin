package service

import (
	"gin/model/table"
	"time"
)

type GroupMsgService struct{}

func (GroupMsgService) SaveGroupMsg(FromUser int, GroupID, Msg string, Attachments string) table.GroupMsg {
	t := table.GroupMsg{FromUser: FromUser, GroupID: GroupID, Msg: Msg, Attachments: Attachments, Time: time.Now()}
	db.Create(&t)
	return t
}
