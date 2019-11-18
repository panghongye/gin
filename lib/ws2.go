package lib

import (
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
		log.Printf("Error: %s, %s\n", err, godebug.LF())
		so.BroadcastTo("", "")
	})

	server.OnAny(func(arg ...interface{}) {
		log.Println(arg)
		log.Println()
	})

	server.On("connection", func(so socketio.Socket) {

		so.On("error", func(msg string) {
			log.Println("[error]", msg)
			log.Println()
		})

		so.On("disconnect", func(msg string) {
			log.Println("【disconnect】", msg)
		})

		so.OnAny(func(arg ...interface{}) {
			log.Println(arg)
			log.Println()
		})
		so.On("test", func() {
			so.Join("")
			so.BroadcastTo("", "")
		})

	})

	return server
}
