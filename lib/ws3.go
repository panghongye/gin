package lib

import (
	"gin/socketio"
	"log"
	"time"
	//  "github.com/zyxar/socketio"
)

func GetWs3() *socketio.Server {
	server, _ := socketio.NewServer(time.Second*25, time.Second*5, socketio.DefaultParser)
	sp := server.Namespace("/")
	sp.OnConnect(func(so socketio.Socket) {
		so.Join("a")
		log.Println("OnConnect <<", so.Sid())
	}).
		OnDisconnect(func(so socketio.Socket) {
			log.Println("OnDisconnect <<", so.Sid())
			so.Close()
		}).
		OnError(func(so socketio.Socket, err ...interface{}) {
			log.Println("OnError <<", so.Sid())
			log.Println(err)
			so.Close()
		}).
		OnEvent("initSocket", func(so socketio.Socket, data uint) {
			log.Println(data)
			so.Emit("initSocket", "xx?")
		}).
		OnEvent("chat message", func(so socketio.Socket, data string) {
			log.Println(data)
			sp.BroadcastToRoom("a", "chat message", data)
		})

	return server
}
