package request

// 分页参数
type PageParam struct {
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"pageSize" binding:"required"`
	MemberID int64
}

// 登录参数
type LoginParam struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Captcha  string `form:"captcha" binding:"required"`
}

// 手机参数
type MobileParam struct {
	Mobile string `form:"mobile" binding:"required"`
}

// 用户名参数
type UserNameParam struct {
	Username string `form:"username" binding:"required"`
}

// 注册参数
type RegParam struct {
	UserName   string `form:"username" binding:"required"`
	Mobile     string `form:"mobile" binding:"required"`
	Captcha    string `form:"captcha" binding:"required"`
	Password   string `form:"password" binding:"required"`
	PayPass    string `form:"paypass" binding:"required"`
	InviteCode string `form:"invitecode" binding:"required"`
	TrueName   string `form:"truename" binding:"required"`
	Alipay     string `form:"alipay" binding:"required"`
	BankName   string `form:"bankname" binding:"required"`
	BankNo     string `form:"bankno" binding:"required"`
}

//登录密码重置参数
type LoginPwdResetParam struct {
	MobileNo      string `form:"mobileNo" binding:"required"`      //手机号码
	Captcha       string `form:"captcha" binding:"required"`       //验证码
	NewPwd        string `form:"newPwd" binding:"required"`        //新密码
	NewPwdConfirm string `form:"newPwdConfirm" binding:"required"` //新密码确认
}
