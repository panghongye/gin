package controller

import (
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

func (this UserCtrl) Login(ctx *gin.Context) {
	res := new(response.Response)
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
		user.Password = ""
		res.Success = true
		token, err := jwt.Jwt.TokenCreate(jwt.PlayLoad{
			"user": user,
		})
		if err != nil {
			res.Success = false
			res.Message = err.Error()
		}
		res.UserInfo = struct {
			table.UserInfo
			Token   string `json:"token"`
			User_id int    `json:"user_id"`
		}{*user, token, int(user.ID)}
	} else {
		res.Success = false
		res.Message = "密码不正确"
	}
	ctx.JSON(http.StatusOK, res)
}
