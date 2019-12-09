package lib

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/kirinlabs/HttpRequest"
)

func GetRandomString(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func NewHttp() *HttpRequest.Request {
	req := HttpRequest.NewRequest()
	// 设置Headers
	req.SetHeaders(map[string]string{
		"Content-Type": "application/x-www-form-urlencoded", //这也是HttpRequest包的默认设置
	})
	// 设置Cookies
	// req.SetCookies(map[string]string{
	// 	"sessionid": "LSIE89SFLKGHHASLC9EETFBVNOPOXNM",
	// })
	return req
}

func Str_Md5(str string) string {
	has := md5.Sum([]byte(str))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
