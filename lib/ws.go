package lib

import (
	"log"

	"github.com/googollee/go-socket.io"
)

func GetWs() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		server.JoinRoom("room1", s)
		log.Println("connected:", s.ID(), s.RemoteAddr())
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
