package main

import (
	"gin/route"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("【当前时间】", time.Now())
	route.BuildRouter().Run(":6666")
}
