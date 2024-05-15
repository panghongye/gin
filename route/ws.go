package route

import (
	"encoding/json"
	"fmt"
	"gin/lib"
	"gin/lib/convert"
	"gin/lib/jwt"
	"gin/lib/redis"
	"gin/model/response"
	"gin/model/table"
	"gin/service"
	"gin/socketio"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	Redis           = redis.Redis
	userService     service.UserService
	groupService    service.GroupService
	groupMsgService service.GroupMsgService
)

type IntList []int

func (s IntList) Len() int           { return len(s) }
func (s IntList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IntList) Less(i, j int) bool { return s[i] < s[j] }

type Param struct {
	Token string `json:"token"`
}

type MessageType string

const (
	tip MessageType = "tip"
	msg MessageType = "msg"
)

type Message struct {
	Type MessageType `json:"type"`
	Msg  string      `json:"msg"`
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

		// 初始化 获取消息
		np.OnEvent("init", func(so socketio.Socket, param Param) response.Response {
			prefix := "【ws init】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			// 每次连接时将 {userId:socketID}  存入 redis
			Redis.Set(fmt.Sprint(userID), so.Sid(), 0)
			type Group struct {
				table.GroupInfo
				Msgs []service.GroupMsg `json:"msgs"`
			}
			groups := []Group{} //群组及消息
			for _, group := range groupService.FindGroupsByUserID(userID) {
				so.Join(group.ID)
				msg := groupService.FindGroupMsgByGroupID(group.ID, 0, 20)
				if group.IsFriend == 1 {
					group.Name = groupService.FindFriendNameByGroupUser(group.ID, userID)
				}
				groups = append(groups, Group{group, msg})
			}
			return response.Response{Data: map[string]interface{}{
				"groups": groups,
			}}
		})

		np.OnEvent("newFriend", func(so socketio.Socket, param struct {
			Param
			GroupID  string `json:"groupID"`
			ToUserID int    `json:"toUserID"`
		}) response.Response {
			prefix := "【newFriend】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			a := IntList{param.ToUserID, userID}
			sort.Sort(a)
			param.GroupID = fmt.Sprint(a[0]) + "," + fmt.Sprint(a[1])
			groupService.CreateGroup(param.GroupID, param.GroupID, param.GroupID, userID, 1)
			groupService.JoinGroup(param.GroupID, userID, param.ToUserID)
			so.Join(param.GroupID)
			sid := Redis.Get(fmt.Sprint(param.ToUserID)).Val()
			// 向对方广播init 重新获取群组列表
			so.EmitTo(sid, "init")
			return response.Response{Data: param}
		})

		np.OnEvent("newGroup", func(so socketio.Socket, param struct {
			Param
			Name    string `json:"name"`
			Intro   string `json:"intro"`
			GroupID string `json:"groupID"`
		}) response.Response {
			prefix := "【newGroup】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			param.GroupID = convert.RandomString(20)
			groupService.CreateGroup(param.Name, param.Intro, param.GroupID, userID, 0)
			groupService.JoinGroup(param.GroupID, userID)
			so.Join(param.GroupID)
			return response.Response{Data: param}
		})

		np.OnEvent("joinGroup", func(so socketio.Socket, param struct {
			Param
			GroupID  string
			UserName string
		}) response.Response {
			prefix := "【ws joinGroup】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			groupService.JoinGroup(param.GroupID, userID)
			so.Join(param.GroupID)
			so.BroadcastToRoom(param.GroupID, "groupMsg", response.Response{Data: Message{tip, "欢迎新的小伙伴：" + param.UserName}})
			return response.Response{}
		})

		np.OnEvent("leaveGroup", func(so socketio.Socket, data struct {
			UserID  string
			GroupID string
		}) {
			so.Leave(data.GroupID)
			groupService.LeaveGroup(data.UserID, data.GroupID)
		})

		np.OnEvent("sendGroupMsg", func(so socketio.Socket, param struct {
			Param
			GroupID     string        `json:"groupID"`
			Time        time.Time     `json:"time"`
			Msg         string        `json:"msg"`
			UserName    string        `json:"userName"`
			Attachments []interface{} `json:"attachments"`
		}) response.Response {
			prefix := "【ws sendGroupMsg】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			msg := struct {
				UserName string `json:"userName"`
				table.GroupMsg
			}{param.UserName, groupMsgService.SaveGroupMsg(userID, param.GroupID, param.Msg, attachmentsTOJsonStr(param.Attachments))}

			so.BroadcastToRoom(param.GroupID, "getGroupMsg", response.Response{Data: msg})
			return response.Response{Data: msg}
		})

		np.OnEvent("updateGroupInfo", func(so socketio.Socket, data struct {
			Name    string `json:"name"`
			Intro   string `json:"intro"`
			GroupID string `json:"groupID"`
		}) string {
			groupService.UpdateGroupInfo(data.Name, data.Intro, data.GroupID)
			return "修改群资料成功"
		})

		np.OnEvent("getGroupInfo", func(so socketio.Socket, param struct {
			Param
			ID string
		}) response.Response {
			prefix := "【ws getGroupInfo】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			t := groupService.FindGroupByID(param.ID)
			if t.ID != "" {
				return response.Response{Data: t}
			}
			return response.Response{}
		})

		np.OnEvent("getUserInfo", func(so socketio.Socket, param struct {
			Param
			ID int
		}) response.Response {
			prefix := "【ws getUserInfo】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			t := new(table.UserInfo)
			userService.FindUserByID(param.ID).Scan(t)
			return response.Response{Data: t}
		})

		np.OnEvent("FindUserByID", func(so socketio.Socket, user_id int) *table.UserInfo {
			t := &table.UserInfo{}
			userService.FindUserByID(user_id).Scan(&t)
			return t
		})

		np.OnEvent("search", func(so socketio.Socket, param struct {
			Param
			Search string
		}) response.Response {
			prefix := "【search】"
			userID := getTokenDataID(prefix, param.Token)
			if userID == 0 {
				return response.Response{Code: response.TokenErr.Code, Msg: response.TokenErr.Msg}
			}
			data := map[string]interface{}{}
			data["users"] = userService.FuzzyFindUsersByName(param.Search)
			data["groups"] = groupService.FuzzyFindGroupsByName(param.Search)
			return response.Response{Data: data}
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

func getTokenDataID(prefix, token string) int {
	data, err := getTokenData(prefix, token)
	if err != nil {
		return 0
	}
	return convert.StringToInt(data.PlayLoad)
}
