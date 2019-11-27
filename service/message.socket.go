package service

import "gin/model/table"

type Message struct{}

type ClientHomePage struct {
	Attachments    string `json:"attachments"`
	Avatar         string `json:"avatar"`
	Be_friend_time int    `json:"be_friend_time"`
	Github_id      string `json:"github_id"`
	Message        string `json:"message"`
	Name           string `json:"name"`
	ShowCallMeTip  bool   `json:"showCallMeTip"`
	Time           int64  `json:"time"`
	Unread         int    `json:"unread"`
	User_id        int    `json:"user_id"`
	To_group_id    string `json:"to_group_id"`
}

var (
	userService      UserService
	chatService      ChatService
	groupChatService GroupChatService
)

func (this Message) GetPrivateMsg(toUser, user_id, start, count int) map[string]interface{} {
	if start <= 0 {
		start = 1
	}
	if start <= 0 {
		count = 20
	}
	messages := []struct {
		Message     string `json:"message"`
		Attachments string `json:"attachments"`
		Time        int    `json:"time"`
		From_user   int    `json:"from_user"`
		To_group_id string `json:"to_group_id"`
		Avatar      string `json:"avatar"`
		Name        string `json:"name"`
		Github_id   string `json:"github_id"`
	}{}
	chatService.GetPrivateDetail(user_id, toUser, start-1, count).Scan(&messages)
	var userInfo table.UserInfo
	userService.GetUserInfo(toUser).Scan(&userInfo)

	return map[string]interface{}{
		"messages": messages,
		"userInfo": userInfo,
	}
}

func (this Message) GetGroupItem(groupId string, start, count int) map[string]interface{} {
	if start <= 0 {
		start = 1
	}
	if start <= 0 {
		count = 20
	}

	messages := []struct {
		Message     string `json:"message"`
		Attachments string `json:"attachments"`
		Time        int    `json:"time"`
		From_user   int    `json:"from_user"`
		To_group_id string `json:"to_group_id"`
		Avatar      string `json:"avatar"`
		Name        string `json:"name"`
		Github_id   string `json:"github_id"`
	}{}

	groupChatService.GetGroupMsg(groupId, start-1, count).Scan(&messages)
	var groupInfo struct {
		To_group_id string `json:"to_group_id"`
		Name        string `json:"name"`
		Creator_id  int    `json:"creator_id"`
		Create_time int    `json:"create_time"`
		Members     []struct {
			User_id   int    `json:"user_id"`
			Socketid  string `json:"socketid"`
			Name      string `json:"name"`
			Avatar    string `json:"avatar"`
			Github_id string `json:"github_id"`
			Github    string `json:"github"`
			Intro     string `json:"intro"`
			Company   string `json:"company"`
			Location  string `json:"location"`
			Website   string `json:"website"`
		} `json:"members"`
	}
	groupChatService.GetGroupInfo(groupId, "").Scan(&groupInfo)
	groupChatService.GetGroupMember(groupId, start-1, count).Scan(&groupInfo.Members)

	return map[string]interface{}{
		"messages":  messages,
		"groupInfo": groupInfo,
	}
}

func (this Message) GetAllMessage(user_id int, clientHomePageList []ClientHomePage) map[string]interface{} {
	privateList := []ClientHomePage{}
	groupList := []ClientHomePage{}
	privateChat := []interface{}{}
	groupChat := []interface{}{}
	userService.GetPrivateList(user_id).Scan(&privateList)
	userService.GetGroupList(user_id).Scan(&groupList)
	homePageList := append(groupList, privateList...)
	for _, item := range homePageList {
		var goal *ClientHomePage
		for _, e := range clientHomePageList {
			var b bool
			if e.User_id == 0 {
				b = e.To_group_id == item.To_group_id
			} else {
				b = e.User_id == item.User_id
			}
			if b {
				goal = &e
			}
		}
		if goal != nil {
			sortTime := goal.Time
			res := ClientHomePage{}
			if item.User_id != 0 {
				chatService.GetUnreadCount(sortTime, user_id, item.User_id).Scan(&res)
			} else {
				groupChatService.GetUnreadCount(sortTime, item.To_group_id).Scan(&res)
			}
			item.Unread = goal.Unread + res.Unread
		}

		if item.User_id != 0 {
			privateChat = append(privateChat, []interface{}{item.User_id, this.GetPrivateMsg(item.User_id, user_id, 0, 0)})
		} else if item.To_group_id != "" {
			groupChat = append(groupChat, []interface{}{item.To_group_id, this.GetGroupItem(item.To_group_id, 0, 0)})
		}
	}

	return map[string]interface{}{
		"homePageList": homePageList,
		"groupChat":    groupChat,
		"privateChat":  privateChat,
	}
}
