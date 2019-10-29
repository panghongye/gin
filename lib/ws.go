package lib

import (
	"fmt"
	"log"
	"time"

	engineio "github.com/googollee/go-engine.io"
	socketio "github.com/googollee/go-socket.io"
)

func GetWs() *socketio.Server {
	server, err := socketio.NewServer(&engineio.Options{nil, nil, time.Hour, 0, nil, nil})
	if err != nil {
		log.Fatal(err)
	}
	server.OnError("/", func(e error) {
		log.Println("error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
		s.Close()
		// server.BroadcastToRoom("room1", "chat message", s.ID()+"离开了")
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		// s.SetContext("")
		// server.JoinRoom("room1", s)
		// log.Println("connected:", s.ID(), s.RemoteAddr())
		// server.BroadcastToRoom("room1", "chat message", s.ID()+"已连接")
		return nil
	})

	server.OnEvent("/", "initSocket", func(s socketio.Conn, msg string) {
		// server.BroadcastToRoom("room1", "chat message", s.ID()+" : "+msg)
		fmt.Println("????", s.Namespace())
		fmt.Println()
	})

	go server.Serve()
	return server
}
