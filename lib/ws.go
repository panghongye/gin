package lib

import (
	"log"
	"time"

	"github.com/googollee/go-engine.io"
	"github.com/googollee/go-socket.io"
)

func GetWs() *socketio.Server {
	server, err := socketio.NewServer(&engineio.Options{nil, nil, time.Hour, time.Second, nil, nil})
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		server.JoinRoom("room1", s)
		log.Println("connected:", s.ID(), s.RemoteAddr())
		server.BroadcastToRoom("room1", "chat message", s.ID()+"已连接")
		return nil
	})

	server.OnError("/", func(e error) {
		log.Println("error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnEvent("/", "chat message", func(s socketio.Conn, msg string) {
		log.Println("chat message:", msg)
		server.BroadcastToRoom("room1", "chat message", msg)
	})

	go server.Serve()
	return server
}
