package middleware

import (
	"errors"
	"gin/lib/jwt"
	"gin/model/response"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Cros(c *gin.Context) {
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Request.Header.Del("Origin")
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
	}
}

func Auth(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		SendErr(errors.New("token 无效"), c)
	}
	_, err := jwt.Jwt.TokenParse(token)
	if err != nil {
		SendErr(err, c)
	}
}

// 有错时 返回 true
func SendErr(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}
	logrus.Info("SendErr ", err)
	c.JSON(200, response.Response{
		Message: err.Error(),
	})
	c.Abort()
	return true
}
