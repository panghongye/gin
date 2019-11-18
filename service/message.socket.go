package service

type Message struct{}

type ClientHomePage struct {
	Ttachments     string `json:"ttachments"`
	Avatar         string `json:"avatar"`
	Be_friend_time int    `json:"be_friend_time"`
	Github_id      string `json:"github_id"`
	Message        string `json:"message"`
	Name           string `json:"name"`
	ShowCallMeTip  bool   `json:"showCallMeTip"`
	Time           int64  `json:"time"`
	Unread         int    `json:"unread"`
	User_id        int    `json:"user_id"`
	To_group_id    string `json:"To_group_id"`
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
	var messages = chatService.GetPrivateDetail(user_id, toUser, start-1, count)
	var userInfo = userService.GetUserInfo(toUser)
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

	var messages = groupChatService.GetGroupMsg(groupId, start-1, count)
	var groupInfo struct {
		Members interface{} `json:"members"`
	}
	groupChatService.GetGroupInfo(groupId, ``).Scan(&groupInfo)
	groupInfo.Members = groupChatService.GetGroupMember(groupId, start-1, count)

	return map[string]interface{}{
		"messages":  messages,
		"groupInfo": groupInfo,
	}
}

func (this Message) GetAllMessage(user_id int, clientHomePageList []ClientHomePage) map[string]interface{} {
	var privateList []ClientHomePage
	var groupList []ClientHomePage
	privateChat := map[interface{}]interface{}{}
	groupChat := map[interface{}]interface{}{}
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
			var res = ClientHomePage{}
			if item.User_id != 0 {
				chatService.GetUnreadCount(sortTime, user_id, item.User_id).Scan(&res)
			} else {
				groupChatService.GetUnreadCount(sortTime, item.To_group_id).Scan(&res)
			}
			item.Unread = goal.Unread + res.Unread
		}

		if item.User_id != 0 {
			// TODO data
			var data interface{} = this.GetPrivateMsg(item.User_id, user_id, 0, 0)
			privateChat[item.User_id] = data
		} else if item.To_group_id != "" {
			var data interface{} = this.GetGroupItem(item.To_group_id, 0, 0)
			groupChat[item.To_group_id] = data
		}

	}

	return map[string]interface{}{
		"homePageList": homePageList,
		"privateChat":  privateChat,
		"groupChat":    groupChat,
	}
}
