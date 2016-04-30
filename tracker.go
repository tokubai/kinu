package main

import (
	"time"
	"strconv"
	"github.com/TakatoshiMaeda/kinu/logger"
	"github.com/Sirupsen/logrus"
	"fmt"
)

func track(tag string, f func()) {
	if !TRACKABLE {
		f()
		return
	}
	now := time.Now()
	f()
	logger.WithFields(logrus.Fields{
		"process_time": strconv.Itoa(int(time.Now().Sub(now) / time.Microsecond)) + "ms",
	}).Debug(fmt.Sprintf("%s process time", tag))
}
