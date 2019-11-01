package lib

import (
	"log"
	"time"

	"github.com/zyxar/socketio"
)

func GetWsx() *socketio.Server {
	server, _ := socketio.NewServer(time.Second*25, time.Second*5, socketio.DefaultParser)
	server.Namespace("/").
		OnConnect(func(so socketio.Socket) {
			log.Println("connected:", so.RemoteAddr(), so.Sid(), so.Namespace())
		}).
		OnDisconnect(func(so socketio.Socket) {
			log.Printf("%v %v %q disconnected", so.Sid(), so.RemoteAddr(), so.Namespace())
			so.Close()
		}).
		OnError(func(so socketio.Socket, err ...interface{}) {
			log.Println("socket", so.Sid(), so.RemoteAddr(), so.Namespace(), "error:", err)
			log.Println(err)
		}).
		OnEvent("initSocket", func(so socketio.Socket, data string) {
			log.Println(data)
		})
	return server
}
