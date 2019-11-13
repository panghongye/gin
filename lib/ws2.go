package lib

import (
	"github.com/pschlump/godebug"
	"github.com/pschlump/socketio"
	"log"
	"strings"
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

		so.On("initSocket", func(userID int) {
			t := userService.GetByID(userID)
			socketId := so.Id()
			if t.Socketid != "" {
				socketId = strings.Split(t.Socketid, ",")[0] + "," + socketId
			}
			if result := userService.SaveUserSocketId(userID, socketId); result.Error != nil {
				so.Emit("error", map[string]interface{}{
					"code": 500, "message": err.Error(),
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
				so.Emit("error", map[string]interface{}{
					"code": 500, "message": err.Error(),
				})
				return
			}
			so.Emit("initSocket success")
		})

		so.On("initMessage", func() {
			so.Emit("initMessage")
		})

		so.On("sendGroupMsg", func() {
			so.Emit("sendGroupMsg")
		})

		so.On("joinGroup", func() {
			so.Emit("joinGroup")
		})

		so.On("sendPrivateMsg", func() {
			so.Emit("sendPrivateMsg")
		})

		so.On("addAsTheContact", func() {
			so.Emit("addAsTheContact")
		})

		so.On("robotChat", func() {
			so.Emit("robotChat")
		})

	})

	return server
}
