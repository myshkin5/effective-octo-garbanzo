package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Logger ExternalLogger
)

type ExternalLogger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

func init() {
	Logger = logrus.WithField("service_name", "effective-octo-garbanzo")
	logrus.SetFormatter(&JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime:  "@timestamp",
			FieldKeyLevel: "priority",
			FieldKeyMsg:   "@message",
		},
		LevelMap: LevelMap{
			logrus.PanicLevel: "PANIC",
			logrus.FatalLevel: "FATAL",
			logrus.ErrorLevel: "ERROR",
			logrus.WarnLevel:  "WARN",
			logrus.InfoLevel:  "INFO",
			logrus.DebugLevel: "DEBUG",
		},
	})
	// Default to no logging for quiet tests
	logrus.SetLevel(logrus.PanicLevel)
}

func Init() error {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	return nil
}
