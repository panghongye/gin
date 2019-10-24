package route

import (
	"doraemon/doraemon-api/controller"
	"doraemon/doraemon-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func BuildRouter() *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	//router.Use(middleware.Recover())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	}))

	// 语言读取
	// router.Use(middleware.Language)

	// 公开组路由
	publicRouter := router.Group("/v1/public")

	//登录（公开）
	loginController := new(controller.LoginController)
	publicRouter.POST("/login/submit", loginController.Login)
	publicRouter.POST("/login/captcha", loginController.LoginCaptcha)
	publicRouter.POST("/login/loginPwdReset", loginController.LoginPwdReset)
	publicRouter.POST("/login/loginPwdResetSendCode", loginController.LoginPwdResetSendCode)
	publicRouter.POST("/login/exist", loginController.Exist)

	//注册（公开）
	registerController := new(controller.RegisterController)
	publicRouter.POST("/reg/submit", registerController.Reg)
	publicRouter.POST("/reg/captcha", registerController.RegCaptcha)

	//用户
	memberRouter := router.Group("/v1/member")
	memberRouter.Use(middleware.Auth)
	memberRouter.Use(middleware.Holiday())
	memberController := new(controller.MemberController)
	memberRouter.POST("/info", memberController.Info)                     //获取用户信息
	memberRouter.POST("/changeLoginPwd", memberController.ChangeLoginPwd) //修改登录密码
	memberRouter.POST("/changePayPwd", memberController.ChangePayPwd)     //修改支付密码
	memberRouter.POST("/sendSmsCode", memberController.SendSmsCode)       //发送短信验证码
	memberRouter.POST("/queryTeamLevel", memberController.QueryTeamLevel) //我的团队-层级查询
	memberRouter.POST("/queryTeamStat", memberController.QueryTeamStat)   //我的团队-统计信息
	memberRouter.POST("/activate", memberController.Activate)             //激活个人帐号

	//首页
	homeController := new(controller.HomeController)
	homeRouter := router.Group("/v1/home")
	homeRouter.Use(middleware.Auth)
	homeRouter.Use(middleware.Holiday())
	homeRouter.POST("/transferType", homeController.TransferType) //转账类型
	homeRouter.POST("/transferInfo", homeController.TransferInfo) //转账信息
	homeRouter.POST("/transfer", homeController.Transfer)         //转账

	//财富
	assetController := new(controller.AssetController)
	assetRouter := router.Group("/v1/asset")
	assetRouter.Use(middleware.Auth)
	assetRouter.Use(middleware.Holiday())
	assetRouter.POST("/staticList", assetController.StaticList)     //静态钱包列表
	assetRouter.POST("/dynamicList", assetController.DynamicList)   //动态钱包列表
	assetRouter.POST("/coinBillList", assetController.CoinBillList) //币种钱包列表
	assetRouter.POST("/luckList", assetController.LuckList)         //幸运单列表
	assetRouter.POST("/assetStat", assetController.AssetStat)       //资产统计

	// 系统
	// 每日数据
	sysController := new(controller.SysDataController)
	sysRouter := router.Group("/v1/sys")
	sysRouter.POST("/sysdata", sysController.SysData)

	// 系统公告
	articleController := new(controller.ArticleController)
	sysRouter.POST("/article/list", articleController.List) //系统公告列表
	sysRouter.POST("/article/info", articleController.Info) //系统公告详情详情

	// banner
	bannerController := new(controller.BannerController)
	sysRouter.POST("/banner/list", bannerController.List) //banner列表

	//订单
	orderController := new(controller.OrderController)
	orderRouter := router.Group("/v1/order")
	orderRouter.Use(middleware.Auth)
	orderRouter.Use(middleware.Holiday())

	orderRouter.POST("/provide", orderController.Provide)                       //提供帮助
	orderRouter.POST("/accept", orderController.Accept)                         //接受帮助
	orderRouter.POST("/acceptOrderList", orderController.AcceptOrderList)       //接受订单列表
	orderRouter.POST("/provideOrderList", orderController.ProvideOrderList)     //提供订单列表
	orderRouter.POST("/acceptDetail", orderController.AcceptDetail)             //接受订单详情
	orderRouter.POST("/provideDetail", orderController.ProvideDetail)           //提供订单详情
	orderRouter.POST("/confirmPay", orderController.ConfirmPay)                 //确认打款
	orderRouter.POST("/confirmReceipt", orderController.ConfirmReceipt)         //确认收款
	orderRouter.POST("/uploadPayImg", orderController.UploadPayImg)             //上传打款图片
	orderRouter.POST("/deletePayImg", orderController.DeletePayImg)             //删除打款图片
	orderRouter.POST("/complaint", orderController.Complaint)                   //投诉
	orderRouter.POST("/complaintInfo", orderController.ComplaintInfo)           //投诉信息
	orderRouter.POST("/uploadComplaintImg", orderController.UploadComplaintImg) //上传投诉图片
	orderRouter.POST("/deleteComplaintImg", orderController.DeleteComplaintImg) //删除投诉图片
	orderRouter.POST("/provideAmountList", orderController.ProvideAmountList)   //提供帮助金额列表
	orderRouter.POST("/acceptWalletList", orderController.AcceptWalletList)     //接受帮助钱包列表

	//抽奖
	drawController := new(controller.DrawController)
	drawRouter := router.Group("/v1/draw")
	drawRouter.POST("/index", drawController.Index)     //抽奖转盘内容
	drawRouter.POST("/list", drawController.List)       //系统中奖列表
	drawRouter.POST("/myDraws", drawController.MyDraws) //我的中奖列表
	drawRouter.POST("/draw", drawController.Draw)       //抽奖动作

	return router
}
