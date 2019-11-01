package lib

import (
	"fmt"
	"log"

	"github.com/pschlump/godebug"
	"github.com/pschlump/socketio"
)

func GetWs2() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("error", func(so socketio.Socket, err error) {
		fmt.Printf("Error: %s, %s\n", err, godebug.LF())
	})

	server.On("connection", func(so socketio.Socket) {

		so.On("disconnect", func(msg string) {
			fmt.Println("【disconnect】", msg)
		})

	})

	return server
}
