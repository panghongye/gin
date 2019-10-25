package main

import (
	"gin/conf"
	"gin/route"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	logrus.Info("【当前时间】", time.Now())
	route.BuildRouter().Run(conf.Port)
}
