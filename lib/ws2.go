package lib

import (
	"fmt"
	"log"

	"strings"

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

		so.On("initSocket", func(userID int) {
			t := userService.GetByID(userID)
			socketId := so.Id()
			if t.Socketid != "" {
				socketId = strings.Split(t.Socketid, ",")[0] + "," + socketId
			}
			if result := userService.SaveUserSocketId(userID, socketId); result.Error != nil {
				so.Emit("error", struct {
					Code    int
					Message string
				}{
					500,
					result.Error.Error(),
				})
				return
			}
			so.Emit("initSocket success")
		})

		so.On("initGroupChat", func(userID int) {
			t := userService.GetByID(userID)
			socketId := so.Id()
			if t.Socketid != "" {
				socketId = strings.Split(t.Socketid, ",")[0] + "," + socketId
			}
			if result := userService.SaveUserSocketId(userID, socketId); result.Error != nil {
				so.Emit("error", struct {
					Code    int
					Message string
				}{
					500,
					result.Error.Error(),
				})
				return
			}
			so.Emit("initSocket success")
		})

		so.On("disconnect", func(msg string) {
			fmt.Println("【disconnect】", msg)
		})
	})

	return server
}
