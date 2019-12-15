package route

import (
	"encoding/json"
	"gin/lib"
	"gin/model/response"
	"gin/model/table"
	"gin/service"
	"gin/socketio"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	userService      service.UserService
	groupService     service.GroupService
	chatService      service.ChatService
	groupChatService service.GroupChatService
	message          service.Message
)

func getWs() *socketio.Server {
	ws, _ := socketio.NewServer(time.Second*25, time.Second*3, socketio.DefaultParser)
	ws.OnError(func(err error) {
		logrus.Error("[ws.OnError]", err)
	})

	{ // 业务
		np := ws.Namespace("/").OnConnect(func(so socketio.Socket) {
			logrus.Info("【连接】<<", so.Sid())
		})

		np.OnError(func(so socketio.Socket, err ...interface{}) {
			logrus.Info("【错误】<<", so.Sid(), err, so.Close())
		})

		np.OnDisconnect(func(so socketio.Socket) {
			logrus.Info("【断开】<<", so.Sid(), so.Close())
		})

		np.OnEvent("init", func(s socketio.Socket, userID int) response.Response {
			var t table.UserInfo
			userService.GetUserInfo(userID).Scan(&t)
			socketId := s.Sid()
			if t.Socketid != "" {
				socketId = strings.Split(t.Socketid, ",")[0] + "," + socketId
			}
			data := response.Response{}
			if result := userService.SaveUserSocketId(userID, socketId); result.Error != nil {
				data = response.Response{Code: 500, Msg: result.Error.Error()}
				return data
			}

			return data
		})

		np.OnEvent("sendGroupMsg", func(so socketio.Socket, data struct {
			From_user   int           `json:"from_user"`
			To_group_id string        `json:"to_group_id"`
			Time        int           `json:"time"`
			Message     string        `json:"message"`
			Name        string        `json:"name"`
			Attachments []interface{} `json:"attachments"`
		}) interface{} {
			data.Time = int(time.Now().Unix())
			groupChatService.SaveGroupMsg(data.From_user, data.To_group_id, data.Message, attachmentsTOJsonStr(data.Attachments))
			so.BroadcastToRoom(data.To_group_id, "getGroupMsg", data)
			return data
		})

		np.OnEvent("getOneGroupMessages", func(so socketio.Socket, data struct {
			Start   int    `json:"start"`
			GroupId string `json:"groupId"`
			Time    int    `json:"time"`
			Count   int    `json:"count"`
		}) interface{} {
			groupMessages := groupChatService.GetGroupMsg(data.GroupId, data.Start-1, data.Count)
			return groupMessages.Value
		})

		np.OnEvent("getOneGroupItem", func(so socketio.Socket, data struct {
			GroupId string `json:"groupId"`
			Start   int    `json:"start"`
			Count   int    `json:"count"`
		}) map[string]interface{} {
			if data.Start < 1 {
				data.Start = 1
			}
			groupMsgAndInfo := message.GetGroupItem(data.GroupId, data.Start, 20)
			return groupMsgAndInfo
		})

		np.OnEvent("newGroup", func(so socketio.Socket, data struct {
			Name         string `json:"name"`
			Group_notice string `json:"group_notice"`
			Creator_id   int    `json:"creator_id"`
			To_group_id  string `json:"to_group_id"`
			Create_time  int    `json:"create_time"`
		}) interface{} {
			data.To_group_id = lib.GetRandomString(90)
			groupService.CreateGroup(data.Name, data.Group_notice, data.To_group_id, data.Creator_id)
			groupService.JoinGroup(data.Creator_id, data.To_group_id)
			so.Join(data.To_group_id)
			return data
		})

		np.OnEvent("updateGroupInfo", func(so socketio.Socket, data struct {
			Name         string `json:"name"`
			Group_notice string `json:"group_notice"`
			To_group_id  string `json:"to_group_id"`
		}) string {
			groupService.UpdateGroupInfo(data.Name, data.Group_notice, data.To_group_id)
			return "修改群资料成功"
		})

		np.OnEvent("joinGroup", func(so socketio.Socket, data struct {
			UserInfo  table.UserInfo
			ToGroupId string
		}) map[string]interface{} {
			t := []interface{}{}
			groupService.IsInGroup(data.UserInfo.ID, data.ToGroupId).Scan(&t)
			if len(t) < 1 {
				groupService.JoinGroup(data.UserInfo.ID, data.ToGroupId)
				so.BroadcastToRoom(data.ToGroupId, "getGroupMsg", struct {
					table.UserInfo
					message     string
					to_group_id string
					tip         string
				}{data.UserInfo, data.UserInfo.Name + "加入了群聊", data.ToGroupId, "joinGroup"})
			}
			so.Join(data.ToGroupId)
			groupItem := message.GetGroupItem(data.ToGroupId, 0, 0)
			return groupItem
		})

		np.OnEvent("leaveGroup", func(so socketio.Socket, data struct {
			User_id   string
			ToGroupId string
		}) {
			so.Leave(data.ToGroupId)
			groupService.LeaveGroup(data.User_id, data.ToGroupId)
		})


		np.OnEvent("getUserInfo", func(so socketio.Socket, user_id int) interface{} {
			t := &struct {
				User_id int `json:"user_id"`
				table.UserInfo
			}{}
			userService.GetUserInfo(user_id).Scan(&t)
			return t
		})

		np.OnEvent("fuzzyMatch", func(so socketio.Socket, data struct {
			Field      string `json:"field"`
			SearchUser bool   `json:"searchUser"`
		}) map[string]interface{} {
			var fuzzyMatchResult interface{}
			data.Field = "%" + data.Field + "%"
			if data.SearchUser {
				fuzzyMatchResult = userService.FuzzyMatchUsers(data.Field)
				logrus.Info()
			} else {
				fuzzyMatchResult = groupService.FuzzyMatchGroups(data.Field)
				logrus.Info()
			}

			return map[string]interface{}{
				"fuzzyMatchResult": fuzzyMatchResult,
				"searchUser":       data.SearchUser,
			}

		})

		np.OnEvent("robotChat", func(so socketio.Socket, data struct {
			User_id   int
			ToGroupId string
			Message   string
		}) map[string]interface{} {
			req := lib.NewHttp()
			resp, err := req.Post("http://www.tuling123.com/openapi/api", map[string]interface{}{
				"key":    "4e348b4a62ca43b5870b16dc58fbcc93",
				"info":   data.Message,
				"userid": data.User_id,
			})
			d, err := resp.Body()
			s := string(d)
			m := map[string]interface{}{}
			err = json.Unmarshal([]byte(s), &m)
			if err != nil {
				logrus.Error(err)
			}
			return m
		})
	}

	{ // test
		np1 := ws.Namespace("/test")

		np1.OnDisconnect(func(so socketio.Socket) {
			so.LeaveAll()
			so.Close()
		})

		np1.OnConnect(func(so socketio.Socket) {
			so.Join("a")
			logrus.Info("连接 <<", so.Sid())
			so.Emit("ack", "foo", func(msg interface{}) {
				logrus.Info(msg) // TODO 可执行？
			})
		})

		np1.OnEvent("chat message", func(so socketio.Socket, data string) string {
			so.BroadcastToRoom("a", "chat message", so.Sid()+":"+data)
			return data
		})
	}
	return ws
}

func attachmentsTOJsonStr(attachments interface{}) string {
	byte, err := json.Marshal(attachments)
	if err != nil {
		return "[]"
	}
	return string(byte)
}
