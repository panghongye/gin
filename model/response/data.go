package response

import "gin/model/table"

type Response struct {
	Code     int32       `json:"code"`
	Message  string      `json:"message"`
	Msg      string      `json:"msg"`
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	UserInfo struct {
		table.UserInfo
		Token   string `json:"token"`
		User_id int   `json:"user_id"`
	} `json:"userInfo"`
}

type PageData struct {
	Page     int32       `json:"page"`
	PageSize int32       `json:"pageSize"`
	Total    int32       `json:"total"`
	Items    interface{} `json:"items"`
}
