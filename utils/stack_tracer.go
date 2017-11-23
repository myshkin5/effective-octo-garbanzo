package utils

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

func InitStackTracer() {
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGQUIT)
		buf := make([]byte, 1<<20)
		for {
			logs.Logger.Info("Stack tracer waiting for SIGQUIT")
			<-sigs
			logs.Logger.Info("=== Received SIGQUIT - goroutine stacks ===")
			stacklen := runtime.Stack(buf, true)
			logs.Logger.Info(string(buf[:stacklen]))
			logs.Logger.Info("=== Stacks traced ===")
		}
	}()
}
