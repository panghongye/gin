package route

import (
	"github.com/googollee/go-socket.io"
	"log"
)

func getWs1()  {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(so socketio.Conn) error {
		log.Println("【连接】<<", )

		return nil
	})

	go server.Serve()
	return server
}
