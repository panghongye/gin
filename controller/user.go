package controller

import (
	"fmt"
	"gin/lib/convert"
	"gin/lib/jwt"
	"gin/model/request"
	"gin/model/response"
	"gin/model/table"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	userService service.UserService
)

type UserCtrl struct{}

func (this UserCtrl) Register(ctx *gin.Context) {
	res := new(response.Response)
	// 绑定参数
	param := new(request.LoginParam)
	if err := ctx.ShouldBind(param); err != nil {
		logrus.Info("注册参数错误: ", err)
		res.Msg = "注册参数错误"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	if t := userService.FindDataByName(param.Name); t.ID != 0 {
		res.Msg = "用户已存在"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	if t := userService.InsertData(&table.UserInfo{Name: param.Name, Password: param.Password}); t.ID == 0 {
		res.Msg = "注册失败"
		res.Success = false
	} else {
		res.Msg = "注册成功"
		res.Success = true
	}
	ctx.JSON(http.StatusOK, res)
}

func (this UserCtrl) Login(ctx *gin.Context) {
	res := new(response.Response)
	var param request.LoginParam
	if err := ctx.ShouldBind(&param); err != nil {
		logrus.Info("登录参数错误: ", err)
		res.Msg = "登录参数错误"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	user := userService.FindDataByName(param.Name)
	if user.ID == 0 {
		res.Msg = "用户不存在"
		res.Success = false
		ctx.JSON(http.StatusOK, res)
		return
	}

	if user.Password == convert.StrToMd5(param.Password) {
		res.Success = true
		token, err := jwt.Singleton.TokenCreate(fmt.Sprint( user.ID))
		if err != nil {
			res.Success = false
			res.Msg = err.Error()
		}
		res.Data = struct {
			table.UserInfo
			Token string `json:"token"`
		}{*user, token}
	} else {
		res.Success = false
		res.Msg = "密码不正确"
	}
	ctx.JSON(http.StatusOK, res)
}
