package lib

import (
	"gin/socketio"
	"log"
	"strings"
	"time"
	//  "github.com/zyxar/socketio"
)

func GetWs3() *socketio.Server {
	server, _ := socketio.NewServer(time.Second*25, time.Second*5, socketio.DefaultParser)
	server.OnError(func(err error) {
		log.Panicln("server.OnError", err)
	})
	sp := server.Namespace("/")
	sp.
		OnError(func(so socketio.Socket, err ...interface{}) {
			log.Println("OnError <<", so.Sid())
			log.Println(err)
			so.Close()
		}).OnDisconnect(func(so socketio.Socket) {
		log.Println("OnDisconnect <<", so.Sid())
		so.Close()
	}).OnEvent("initSocket", func(s socketio.Socket, userID int) {
		t := userService.GetByID(userID)
		socketId := s.Sid()
		if t.Socketid != "" {
			socketId = strings.Split(t.Socketid, ",")[0] + "," + socketId
		}
		if result := userService.SaveUserSocketId(userID, socketId); result.Error != nil {
			s.Emit("error", struct {
				Code    int
				Message string
			}{
				500,
				result.Error.Error(),
			})
			return
		}
		s.Emit("initSocket success", "XXX")
	}).OnEvent("initGroupChat", func(s socketio.Socket, userID int) {
		s.Emit("initGroupChat success")
	})

	// assets
	{
		sp := server.Namespace("assets")
		sp.OnConnect(func(so socketio.Socket) {
			so.Join("a")
			log.Println("OnConnect <<", so.Sid())
		}).
			OnEvent("chat message", func(so socketio.Socket, data string) {
				log.Println("chat message:", data)
				sp.BroadcastToRoom("a", "chat message", so.Sid()+":"+data)
			})
	}

	return server
}
