package service

import "github.com/jinzhu/gorm"

type ChatService struct{}

func (ChatService) GetPrivateDetail(from_user, to_user, start, count int) *gorm.DB {
	return db.Raw(`SELECT * FROM ( SELECT p.from_user,p.to_user,p.message,p.attachments,p.time,i.avatar,i.name, i.github_id from private_msg as p inner join user_info as i  on p.from_user = i.id  where  (p.from_user = ? AND p.to_user   = ? )  or (p.from_user = ? AND p.to_user   = ? )  order by time desc limit ?,? ) as n order by n.time`, from_user, to_user, start, count)
}

func (ChatService) SavePrivateMsg(from_user, to_user int, time int64, message, attachments string) *gorm.DB {
	sql := `INSERT INTO private_msg(from_user,to_user,message,time,attachments)  VALUES(?,?,?,?,?);`
	return db.Raw(sql, from_user, to_user, message, time, attachments)
}

func (ChatService) GetUnreadCount(sortTime int64, from_user int, to_user int) *gorm.DB {
	return db.Raw(`SELECT count(time) as unread FROM private_msg AS p WHERE p.time > ? and ((p.from_user = ? and p.to_user= ?) or (p.from_user = ? and p.to_user=?));`, sortTime, from_user, to_user, to_user, from_user)
}
