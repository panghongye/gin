package request

type PageParam struct {
	Page     int32 `form:"page" binding:"required"`
	PageSize int32 `form:"pageSize" binding:"required"`
	MemberID int64
}

type LoginParam struct {
	Name     string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
	// Captcha  string `form:"captcha" binding:"required"`
}
