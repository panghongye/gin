package main

import (
	_ "gin/conf"
	"gin/route"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	port := viper.GetString("server.addr")
	logrus.Info("[start]" + port)
	if err := route.BuildRouter().Run(port); err != nil {
		panic(err)
	}
}
