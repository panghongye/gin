package lib

import (
	"gin/model/table"
	"gin/service"
	"gin/socketio"
	"log"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	//  "github.com/zyxar/socketio"
)

var (
	userService      service.UserService
	groupService     service.GroupService
	chatService      service.ChatService
	groupChatService service.GroupChatService
	message          service.Message
)

func GetWs() *socketio.Server {
	server, _ := socketio.NewServer(time.Second*25, time.Second*5, socketio.DefaultParser)
	server.OnError(func(err error) {
		logrus.Error("[server.OnError]", err)
	})
	server.Namespace("/").
		OnError(func(so socketio.Socket, err ...interface{}) {
			log.Println("OnError <<", so.Sid())
			logrus.Error("[so.OnError]", err)
		}).
		OnDisconnect(func(so socketio.Socket) {
			log.Println("OnDisconnect:<<", so.Sid())
			so.Close()
		}).
		OnEvent("initSocket", func(s socketio.Socket, userID int) {
			var t table.UserInfo
			userService.GetUserInfo(userID).Scan(&t)
			socketId := s.Sid()
			if t.Socketid != "" {
				socketId = strings.Split(t.Socketid, ",")[0] + "," + socketId
			}
			if result := userService.SaveUserSocketId(userID, socketId); result.Error != nil {
				s.Emit("error", struct {
					Code    int
					Message string
				}{
					500,
					result.Error.Error(),
				})
				return
			}
			s.Emit("initSocket success")
		}).
		OnEvent("initGroupChat", func(so socketio.Socket, userID int) {
			t := userService.GetGroupList(userID)
			// for item := range t.Value {
			// 	so.Join(item.to_group_id)
			// }
			so.Emit("initGroupChat success", t)
		}).
		// 获取群聊和私聊的数据
		OnEvent("initMessage", func(so socketio.Socket, obj struct {
			User_id            int
			ClientHomePageList []service.ClientHomePage
		}) {
			t := message.GetAllMessage(obj.User_id, obj.ClientHomePageList)
			so.Emit("initGroupChat success", t)
		}).
		// sendPrivateMsg
		OnEvent("sendPrivateMsg", func(so socketio.Socket, data struct {
			From_user   int    `json:"from_user"`
			To_user     int    `json:"to_user"`
			Time        int64  `json:"time"`
			Message     string `json:"message"`
			Attachments string `json:"attachments"`
		}) {
			data.Time = time.Now().Unix()
			chatService.SavePrivateMsg(data.From_user, data.To_user, data.Time, data.Message, data.Attachments)
			var t struct {
				Socketid string
			}
			userService.GetUserSocketId(data.To_user).Scan(&t)
			existSocketIdStr := t.Socketid
			if existSocketIdStr != "" {
				toUserSocketIds := strings.Split(existSocketIdStr, ",")
				for _, e := range toUserSocketIds {
					so.BroadcastToRoom(e, "getPrivateMsg", data)
				}
			}
			so.Emit("sendPrivateMsg", data)
		}).
		// 群聊发信息
		OnEvent("sendGroupMsg", func(so socketio.Socket, data struct {
			From_user   int    `json:"from_user"`
			To_group_id string `json:"to_group_id"`
			Time        int64  `json:"time"`
			Message     string `json:"message"`
			Attachments string `json:"attachments"`
		}) {
			data.Time = time.Now().Unix()
			groupChatService.SaveGroupMsg(data.From_user, data.To_group_id, data.Time, data.Message, data.Attachments)
			so.BroadcastToRoom(data.To_group_id, "getGroupMsg", data)
			so.Emit("sendGroupMsg", data)
		}).
		OnEvent("getOnePrivateChatMessages", func(so socketio.Socket, data struct {
			User_id int `json:"user_id"`
			ToUser  int `json:"toUser"`
			Start   int `json:"start"`
			Count   int `json:"count"`
		}) {
			privateMessages := chatService.GetPrivateDetail(data.User_id, data.ToUser, data.Start-1, data.Count)
			so.Emit("getOnePrivateChatMessages", privateMessages)
		}).
		// get group messages in a group;
		OnEvent("getOneGroupMessages", func(so socketio.Socket, data struct {
			Start   int    `json:"start"`
			GroupId string `json:"groupId"`
			Time    int64  `json:"time"`
			Count   int    `json:"count"`
		}) {
			groupMessages := groupChatService.GetGroupMsg(data.GroupId, data.Start-1, data.Count)
			so.Emit("getOneGroupMessages", groupMessages)
		}).
		// get group item including messages and group info.
		OnEvent("getOneGroupItem", func(so socketio.Socket, data struct {
			GroupId string `json:"groupId"`
			Start   int    `json:"start"`
			Count   int    `json:"count"`
		}) {
			if data.Start < 1 {
				data.Start = 1
			}
			groupMsgAndInfo := message.GetGroupItem(data.GroupId, data.Start, 20)
			so.Emit("getOneGroupItem", groupMsgAndInfo)
		}).
		OnEvent("createGroup", func(so socketio.Socket, data struct {
			Name         string `json:"name"`
			Group_notice string `json:"group_notice"`
			Creator_id   uint   `json:"creator_id"`
			To_group_id  string `json:"to_group_id"`
			Create_time  int64  `json:"create_time"`
		}) {
			data.Create_time = time.Now().Unix()
			data.To_group_id = GetRandomString(90)
			groupService.CreateGroup(data.Name, data.Group_notice, data.To_group_id, data.Creator_id, data.Create_time)
			groupService.JoinGroup(data.Creator_id, data.To_group_id)
			so.Join(data.To_group_id)
			so.Emit("createGroup", data)
		})

	// assets
	{
		server.Namespace("/test").
			OnConnect(func(so socketio.Socket) {
				so.Join("a")
				so.Namespace()
				log.Println("OnConnect <<", so.Sid())
			}).
			OnEvent("chat message", func(so socketio.Socket, data string) {
				log.Println("chat message:", data)
				so.BroadcastToRoom("a", "chat message", so.Sid()+":"+data)
			})
	}

	return server
}
