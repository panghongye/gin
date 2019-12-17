package response

// 响应码
type ResCode struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

var (
	Success = ResCode{ 0,  "操作成功"}
	TokenErr       = ResCode{ 1,  "token err"}
	ParamErr       = ResCode{ 1000,  "参数错误"}
)
