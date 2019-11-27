package syserr

type CommonErr struct {
	Label string
	Log   string
	Code  int32       `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//参数错误
type ParamErr struct {
	Label string
	Log   string
}
