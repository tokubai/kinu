package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	ltsv "github.com/doloopwhile/logrusltsv"
	"os"
	"runtime"
	"strconv"
	"time"
)

func init() {
	logLevel := os.Getenv("KINU_LOG_LEVEL")

	if len(logLevel) == 0 {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Panic(err)
	}
	logrus.SetLevel(level)

	formatterType := os.Getenv("KINU_LOG_FORMAT")
	switch formatterType {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "ltsv":
		logrus.SetFormatter(&ltsv.Formatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logrus.SetOutput(os.Stdout)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func Debug(args ...interface{}) {
	logrus.Debug(args)
}

func Info(args ...interface{}) {
	logrus.Info(args)
}

func Warn(args ...interface{}) {
	logrus.Warn(args)
}

func Error(args ...interface{}) {
	logrus.Error(args)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args)
}

func Panic(args ...interface{}) {
	logrus.Panic(args)
}

func ErrorDebug(err error) error {
	if logrus.GetLevel() > logrus.DebugLevel {
		return err
	}

	_, file, line, _ := runtime.Caller(1)
	WithFields(logrus.Fields{
		"file": file + ":" + strconv.Itoa(line),
	}).Error(err.Error())

	return err
}

func TrackResult(tag string, startTime time.Time) {
	WithFields(logrus.Fields{
		"process_time": strconv.Itoa(int(time.Now().Sub(startTime)/time.Millisecond)) + "ms",
	}).Debug(fmt.Sprintf("%s process time", tag))
}
