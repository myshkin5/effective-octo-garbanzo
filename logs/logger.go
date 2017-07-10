package logs

import (
	"os"

	"github.com/op/go-logging"
)

var (
	Logger   *logging.Logger
	LogLevel logging.LeveledBackend
)

func init() {
	Logger = logging.MustGetLogger("go-toccata")
	format := logging.MustStringFormatter("%{time:2006-01-02T15:04:05.000000Z} %{level} %{message}")
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	LogLevel = logging.AddModuleLevel(backendFormatter)
	logging.SetBackend(LogLevel)
	// Default to no logging for quiet tests
	LogLevel.SetLevel(-1, "")
}

func Init(logLevel string) error {
	level, err := logging.LogLevel(logLevel)
	if err != nil {
		return err
	}

	LogLevel.SetLevel(level, "")

	return nil
}
