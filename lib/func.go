package lib

import (
	"log"

	"github.com/gin-gonic/gin"
)

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
