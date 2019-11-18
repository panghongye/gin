package response

// 响应码
type ResCode struct {
	Code int32
	Msg  string
}

var (
	// 系统
	ResCodeSuccess              = &ResCode{Code: 0, Msg: "操作成功"}
	ResCodeFail                 = &ResCode{Code: 200, Msg: "操作失败"}
	ResCodeServerError          = &ResCode{Code: 500, Msg: "服务器错误"}
	ResCodeParamError           = &ResCode{Code: 501, Msg: "参数错误"}
	ResCodeLoginDisabled        = &ResCode{Code: 502, Msg: "系统暂停登录"}
	ResCodeRegisterDisabled     = &ResCode{Code: 503, Msg: "系统暂停注册"}
	ResCodeGetLockKeyError      = &ResCode{Code: 601, Msg: "获取分布式锁异常"}
	ResCodeCommitError          = &ResCode{Code: 701, Msg: "提交异常"}
	ResCodeRequestTooFrequently = &ResCode{Code: 702, Msg: "请求过于频繁"}
)
