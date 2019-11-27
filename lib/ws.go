package lib

import (
	"encoding/json"
	"gin/model/table"
	"gin/service"
	"gin/socketio"
	"log"
	"strings"
	"time"

	"github.com/kirinlabs/HttpRequest"
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
		OnEvent("initSocket", func(s socketio.Socket, userID int) string {
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
			}
			return "initSocket success"
		}).
		OnEvent("initGroupChat", func(so socketio.Socket, userID int) string {
			t := []table.Group_info{}
			userService.GetGroupList(userID).Scan(&t)
			for _, item := range t {
				so.Join(item.To_group_id)
			}
			return "init group chat success"
		}).
		OnEvent("initMessage", func(so socketio.Socket, obj struct {
			User_id            int
			ClientHomePageList []service.ClientHomePage
		}) map[string]interface{} {
			t := message.GetAllMessage(obj.User_id, obj.ClientHomePageList)
			return t
		}).
		OnEvent("sendPrivateMsg", func(so socketio.Socket, data struct {
			From_user   int    `json:"from_user"`
			To_user     int    `json:"to_user"`
			Time        int64  `json:"time"`
			Message     string `json:"message"`
			Attachments string `json:"attachments"`
		}) interface{} {
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
			return data
		}).
		OnEvent("sendGroupMsg", func(so socketio.Socket, data struct {
			From_user   int    `json:"from_user"`
			To_group_id string `json:"to_group_id"`
			Time        int64  `json:"time"`
			Message     string `json:"message"`
			// Attachments []interface{} `json:"attachments"`
		}) interface{} {
			data.Time = time.Now().Unix()
			groupChatService.SaveGroupMsg(data.From_user, data.To_group_id, data.Time, data.Message, "")
			so.BroadcastToRoom(data.To_group_id, "getGroupMsg", data)
			return data
		}).
		OnEvent("getOnePrivateChatMessages", func(so socketio.Socket, data struct {
			User_id int `json:"user_id"`
			ToUser  int `json:"toUser"`
			Start   int `json:"start"`
			Count   int `json:"count"`
		}) interface{} {
			privateMessages := chatService.GetPrivateDetail(data.User_id, data.ToUser, data.Start-1, data.Count)
			return privateMessages.Value
		}).
		// get group messages in a group;
		OnEvent("getOneGroupMessages", func(so socketio.Socket, data struct {
			Start   int    `json:"start"`
			GroupId string `json:"groupId"`
			Time    int64  `json:"time"`
			Count   int    `json:"count"`
		}) interface{} {
			groupMessages := groupChatService.GetGroupMsg(data.GroupId, data.Start-1, data.Count)
			return groupMessages.Value
		}).
		// get group item including messages and group info.
		OnEvent("getOneGroupItem", func(so socketio.Socket, data struct {
			GroupId string `json:"groupId"`
			Start   int    `json:"start"`
			Count   int    `json:"count"`
		}) map[string]interface{} {
			if data.Start < 1 {
				data.Start = 1
			}
			groupMsgAndInfo := message.GetGroupItem(data.GroupId, data.Start, 20)
			return groupMsgAndInfo
		}).
		OnEvent("createGroup", func(so socketio.Socket, data struct {
			Name         string `json:"name"`
			Group_notice string `json:"group_notice"`
			Creator_id   int    `json:"creator_id"`
			To_group_id  string `json:"to_group_id"`
			Create_time  int64  `json:"create_time"`
		}) interface{} {
			data.Create_time = time.Now().Unix()
			data.To_group_id = GetRandomString(90)
			groupService.CreateGroup(data.Name, data.Group_notice, data.To_group_id, data.Creator_id, data.Create_time)
			groupService.JoinGroup(data.Creator_id, data.To_group_id)
			so.Join(data.To_group_id)
			return data
		}).
		OnEvent("updateGroupInfo", func(so socketio.Socket, data struct {
			Name         string `json:"name"`
			Group_notice string `json:"group_notice"`
			To_group_id  string `json:"to_group_id"`
		}) string {
			groupService.UpdateGroupInfo(data.Name, data.Group_notice, data.To_group_id)
			return "修改群资料成功"
		}).
		OnEvent("joinGroup", func(so socketio.Socket, data struct {
			UserInfo  map[string]interface{}
			ToGroupId string
		}) map[string]interface{} {
			if groupService.IsInGroup(data.UserInfo["user_id"].(int), data.ToGroupId).RowsAffected < 1 {
				groupService.JoinGroup(data.UserInfo["user_id"].(int), data.ToGroupId)
				d := map[string]interface{}{
					"message":     data.UserInfo["name"].(string) + "加入了群聊",
					"to_group_id": data.ToGroupId,
					"tip":         "joinGroup",
				}
				for k, v := range data.UserInfo {
					d[k] = v
				}
				so.BroadcastToRoom(data.ToGroupId, "getGroupMsg", d)
			}
			so.Join(data.ToGroupId)
			groupItem := message.GetGroupItem(data.ToGroupId, 0, 0)
			return groupItem
		}).
		OnEvent("leaveGroup", func(so socketio.Socket, data struct {
			User_id   string
			ToGroupId string
		}) {
			so.Leave(data.ToGroupId)
			groupService.LeaveGroup(data.User_id, data.ToGroupId)
		}).
		OnEvent("getGroupMember", func(so socketio.Socket, data struct {
			User_id   string
			ToGroupId string
		}) {
			groupChatService.GetGroupMember(data.ToGroupId, 0, 0)
		}).
		OnEvent("fuzzyMatch", func(so socketio.Socket, data struct {
			User_id   string
			ToGroupId string
		}) {
		}).
		OnEvent("robotChat", func(so socketio.Socket, data struct {
			User_id   int
			ToGroupId string
			Message   string
		}) map[string]interface{} {
			req := NewHttp()
			resp, err := req.Post("http://www.tuling123.com/openapi/api", map[string]interface{}{
				"key":    "4e348b4a62ca43b5870b16dc58fbcc93",
				"info":   data.Message,
				"userid": data.User_id,
			})
			d, err := resp.Body()
			s := string(d)
			m := map[string]interface{}{}
			err = json.Unmarshal([]byte(s), &m)
			log.Println(err)
			return m
		})

	// assets
	{
		server.Namespace("/test").
			OnDisconnect(func(so socketio.Socket) {
				so.LeaveAll()
				so.Close()
			}).
			OnConnect(func(so socketio.Socket) {
				so.Join("a")
				log.Println("OnConnect <<", so.Sid())
			}).
			OnEvent("chat message", func(so socketio.Socket, data string) string {
				so.BroadcastToRoom("a", "chat message", so.Sid()+":"+data)
				return data
			})
	}

	return server
}

func NewHttp() *HttpRequest.Request {
	req := HttpRequest.NewRequest()
	// 设置Headers
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})
	// 设置Cookies
	// req.SetCookies(map[string]string{
	// 	"sessionid": "LSIE89SFLKGHHASLC9EETFBVNOPOXNM",
	// })
	return req
}
