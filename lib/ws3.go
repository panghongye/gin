package lib

import (
	"log"
	"time"

	engineio "github.com/googollee/go-engine.io"
	socketio "github.com/googollee/go-socket.io"
	// socketio "github.com/mrfoe7/go-socket.io"
)

func GetWs3() *socketio.Server {
	server, err := socketio.NewServer(&engineio.Options{nil, nil, time.Hour, time.Second, nil, nil})
	if err != nil {
		log.Fatal(err)
	}

	server.OnEvent("/", "initGroupChat", func(s socketio.Conn, userID int) {
		s.Emit("initGroupChat success")
		// server.BroadcastToRoom("")
	})

	server.OnEvent("/", "test", func(s socketio.Conn, userID int) {
		s.Join("")
		server.BroadcastToRoom("", "")
	})

	go server.Serve()
	return server
}
