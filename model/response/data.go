package response

// 响应码
type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 分页数据
type PageData struct {
	Page     int32       `json:"page"`
	PageSize int32       `json:"pageSize"`
	Total    int32       `json:"total"`
	Items    interface{} `json:"items"`
}
