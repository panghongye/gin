package response

type Response struct {
	Code    int32       `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type PageData struct {
	Page     int32       `json:"page"`
	PageSize int32       `json:"pageSize"`
	Total    int32       `json:"total"`
	Items    interface{} `json:"items"`
}
