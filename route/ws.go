package route

import (
	"encoding/json"
	"gin/lib"
	"gin/lib/convert"
	"gin/lib/jwt"
	"gin/model/response"
	"gin/model/table"
	"gin/service"
	"gin/socketio"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	userService  service.UserService
	groupService service.GroupService
)

type Param struct {
	Token string `json:"token"`
}

func getWs() *socketio.Server {
	ws, _ := socketio.NewServer(time.Second*25, time.Second*3, socketio.DefaultParser)
	ws.OnError(func(err error) {
		logrus.Error("[ws.OnError]", err)
	})

	{ // 业务
		np := ws.Namespace("/").OnConnect(func(so socketio.Socket) {
			prefix := "【连接】"
			logrus.Info(prefix+"<<", so.Sid())
			// fmt.Println("已连接数")
		})

		np.OnError(func(so socketio.Socket, err ...interface{}) {
			logrus.Info("【错误】<<", so.Sid(), err, so.Close())
		})

		np.OnDisconnect(func(so socketio.Socket) {
			logrus.Info("【断开】<<", so.Sid(), so.Close())
		})

		np.OnEvent("init", func(so socketio.Socket, param Param) response.Response {
			prefix := "【ws init】"
			token, err := getTokenData(prefix, param.Token)
			if err != nil {
				return response.Response{Msg: err.Error(), Success: false}
			}

			var t table.UserInfo
			userService.GetUserInfo(convert.StringToInt(token.PlayLoad)).Scan(&t)
			return response.Response{Success: true, Data: t}
			// todo 获取消息列表

		})

		np.OnEvent("sendGroupMsg", func(so socketio.Socket, data struct {
			FromUser    int           `json:"from_user"`
			GroupID     string        `json:"groupId"`
			Time        time.Time     `json:"time"`
			Message     string        `json:"message"`
			Name        string        `json:"name"`
			Attachments []interface{} `json:"attachments"`
		}) interface{} {
			data.Time = time.Now()
			// groupChatService.SaveGroupMsg(data.FromUser, data.GroupID, data.Message, attachmentsTOJsonStr(data.Attachments))
			so.BroadcastToRoom(data.GroupID, "getGroupMsg", data)
			return data
		})

		np.OnEvent("newGroup", func(so socketio.Socket, data struct {
			Name        string `json:"name"`
			GroupNotice string `json:"groupNotice"`
			UserId      int    `json:"userID"`
			GroupID     string `json:"groupId"`
			CreateTime  int    `json:"createTime"`
		}) interface{} {
			data.GroupID = convert.RandomString(20)
			groupService.CreateGroup(data.Name, data.GroupNotice, data.GroupID, data.UserId)
			groupService.JoinGroup(data.UserId, data.GroupID)
			so.Join(data.GroupID)
			return data
		})

		np.OnEvent("newContact", func(so socketio.Socket, data struct {
			UserId int `json:"userID"`
		}) interface{} {
			// data.GroupID = lib.GetRandomString(90)
			// groupService.CreateGroup(data.Name, data.GroupNotice, data.GroupID, data.UserId)
			// groupService.JoinGroup(data.UserId, data.GroupID)
			// so.Join(data.GroupID)
			return data
		})

		np.OnEvent("updateGroupInfo", func(so socketio.Socket, data struct {
			Name        string `json:"name"`
			GroupNotice string `json:"group_notice"`
			GroupID     string `json:"groupId"`
		}) string {
			groupService.UpdateGroupInfo(data.Name, data.GroupNotice, data.GroupID)
			return "修改群资料成功"
		})

		np.OnEvent("joinGroup", func(so socketio.Socket, data struct {
			UserId  int
			GroupID string
		}) {
		})

		np.OnEvent("leaveGroup", func(so socketio.Socket, data struct {
			UserID  string
			GroupID string
		}) {
			so.Leave(data.GroupID)
			groupService.LeaveGroup(data.UserID, data.GroupID)
		})

		np.OnEvent("getUserInfo", func(so socketio.Socket, user_id int) *table.UserInfo {
			t := &table.UserInfo{}
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
			UserID    int
			ToGroupId string
			Message   string
		}) map[string]interface{} {
			req := lib.NewHttp()
			resp, err := req.Post("http://www.tuling123.com/openapi/api", map[string]interface{}{
				"key":    "4e348b4a62ca43b5870b16dc58fbcc93",
				"info":   data.Message,
				"userid": data.UserID,
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

func getTokenData(prefix, token string) (*jwt.Claims, error) {
	_prefix := "getTokenData "
	data, err := jwt.Singleton.TokenParse(token)
	if err != nil {
		logrus.Errorln(prefix+_prefix+"失败", err, token)
	} else {
		logrus.Infoln(prefix+_prefix+"成功", token)
	}
	return data, err
}
