package lib

import (
	"gin/service"
	"log"
	"strings"

	// engineio "github.com/googollee/go-engine.io"
	socketio "github.com/googollee/go-socket.io"
)

var userService service.UserService

func GetWs() *socketio.Server {
	// server, err := socketio.NewServer(&engineio.Options{nil, nil, time.Hour, 0, nil, nil})
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnError("/", func(e error) {
		log.Println("?????????????????:", e.Error())
		log.Println()
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		s.Close()
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		return nil
	})

	server.OnEvent("/", "initSocket", func(s socketio.Conn, userID int) {
		// server.BroadcastToRoom("room1", "chat message", s.ID()+" : "+msg)
		t := userService.GetByID(userID)
		socketId := s.ID()
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
		s.Emit("initSocket success")
	})

	server.OnEvent("/", "initGroupChat", func(s socketio.Conn, userID int) {
		// server.BroadcastToRoom("room1", "chat message", s.ID()+" : "+msg)
		t := userService.GetByID(userID)
		socketId := s.ID()
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
		s.Emit("initSocket success")
	})

	go server.Serve()
	return server
}
