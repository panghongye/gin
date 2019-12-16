package lib

import (
	"github.com/kirinlabs/HttpRequest"
)

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
