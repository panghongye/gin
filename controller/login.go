package controller

import (
	"gin/model/request"
	"gin/model/response"
	"gin/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Login struct{}

//登录
func (c *Login) Login(ctx *gin.Context) {
	res := new(response.Response)

	// 绑定参数
	var param request.LoginParam
	if err := ctx.ShouldBind(&param); err != nil {

		res.Code = response.ResCodeParamError.Code
		res.Msg = response.ResCodeParamError.Msg
		logrus.Info("登录参数错误: ", err)
		ctx.JSON(http.StatusOK, res)
		return
	}

	res = service.UserInfo.Login(&param)

	ctx.JSON(http.StatusOK, res)
}

//检查用户是否已存在
func (c *Login) Exist(ctx *gin.Context) {

	prefix := "【检查用户是否已存在】"
	res := new(response.Response)
	res.Code = response.ResCodeSuccess.Code
	res.Msg = response.ResCodeSuccess.Msg

	// 绑定参数
	var param request.ExistParam
	if err := ctx.ShouldBind(&param); err != nil {

		res.Code = response.ResCodeParamError.Code
		res.Msg = response.ResCodeParamError.Msg

		logrus.Info(prefix, "参数错误: ", err)

		ctx.JSON(http.StatusOK, res)
		return
	}

	res = memberService.Exist(param.Account)

	ctx.JSON(http.StatusOK, res)
}
