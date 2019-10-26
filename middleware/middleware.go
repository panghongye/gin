package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
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

// 有错时 返回 true
func SendErr(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}
	log.Println("SendErr ", err)
	c.JSON(400, gin.H{
		"err": err.Error(),
	})
	return true
}
