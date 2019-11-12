package lib

import (
	"gin/socketio"
	"log"
	"time"
	// "github.com/zyxar/socketioocketio"
)

func GetWs3() *socketio.Server {
	server, _ := socketio.NewServer(time.Second*25, time.Second*5, socketio.DefaultParser)
	server.Namespace("/").
		OnConnect(func(so socketio.Socket) {
			log.Println("connected:", so.RemoteAddr(), so.LocalAddr(), so.Sid(), so.Namespace())
		}).
		OnDisconnect(func(so socketio.Socket) {
			log.Printf("%v %v %q disconnected", so.Sid(), so.RemoteAddr(), so.Namespace())
			so.Close()
		}).
		OnError(func(so socketio.Socket, err ...interface{}) {
			log.Println("socket", so.Sid(), so.RemoteAddr(), so.Namespace(), "error:", err)
			log.Println(err)
		}).
		OnEvent("initSocket", func(so socketio.Socket, data uint) {
			log.Println(data)
			so.Emit("initSocket", "xx?")
		}).
		OnEvent("chat message", func(so socketio.Socket, data string) {
			log.Println(data)
			so.Emit("chat message", data+">>")
		})

	return server
}
