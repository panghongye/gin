package controller

import (
	"gin/model/request"
	"gin/model/response"
	"gin/model/table"
	"gin/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	userService service.UserService
)

type UserCtrl struct{}

// 注册
func (this UserCtrl) Register(ctx *gin.Context) {
	res := new(response.Response)
	// 绑定参数
	param := new(request.LoginParam)
	if err := ctx.ShouldBind(param); err != nil {
		logrus.Info("注册参数错误: ", err)
		res.Message = "注册参数错误"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	if t := userService.FindDataByName(param.Name); t.ID != 0 {
		res.Message = "用户已存在"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	if t := userService.InsertData(&table.UserInfo{Name: param.Name, Password: param.Password}); t.ID == 0 {
		res.Message = "注册失败"
		res.Success = false
	} else {
		res.Message = "注册成功"
		res.Success = true
	}
	ctx.JSON(http.StatusOK, res)
}

//登录
func (this UserCtrl) Login(ctx *gin.Context) {
	res := new(response.Response)
	// 绑定参数
	var param request.LoginParam
	if err := ctx.ShouldBind(&param); err != nil {
		logrus.Info("登录参数错误: ", err)
		res.Message = "登录参数错误"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	user := userService.FindDataByName(param.Name)
	if user.ID == 0 {
		res.Message = "用户不存在"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	if user.Password == param.Password {
		res.Success = true
		user.Password = ""
		res.UserInfo = user
	} else {
		res.Success = false
		res.Message = "密码不正确"
	}
	ctx.JSON(http.StatusOK, res)
}
