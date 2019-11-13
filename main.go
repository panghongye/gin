package main

import (
	_ "gin/conf"
	"gin/route"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.Info("[start]")
	route.BuildRouter().Run(viper.GetString("server.addr"))
}
