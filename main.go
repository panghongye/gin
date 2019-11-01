package main

import (
	"gin/conf"
	"gin/route"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("【启动】")
	route.BuildRouter().Run(conf.Port)
}
