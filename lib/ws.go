package lib

import (
	"gin/service"
	"log"
	"strings"
	"time"

	engineio "github.com/googollee/go-engine.io"
	socketio "github.com/googollee/go-socket.io"
	// socketio "github.com/mrfoe7/go-socket.io"
)

var userService service.UserService

func GetWs() *socketio.Server {
	server, err := socketio.NewServer(&engineio.Options{nil, nil, time.Hour, time.Second, nil, nil})
	if err != nil {
		log.Fatal(err)
	}

	server.OnEvent("/", "initSocket", func(s socketio.Conn, userID int) {
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
		s.Emit("initGroupChat success")
	})

	go server.Serve()
	return server
}
