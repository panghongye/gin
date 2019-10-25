package request

// 分页参数
type PageParam struct {
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"pageSize" binding:"required"`
	MemberID int64
}

// 注册登录参数
type LoginParam struct {
	Name     string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
	// Captcha  string `form:"captcha" binding:"required"`
}
