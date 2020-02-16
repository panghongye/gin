package main

import (
	"gin/route"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.Info("[start]")
	if err := route.BuildRouter().Run(viper.GetString("server.addr")); err != nil {
		panic(err)
	}
}
