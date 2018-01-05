package handlers

import (
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

type loggingWriter struct{}

func LoggingHandler(handler http.Handler) http.Handler {
	return handlers.LoggingHandler(loggingWriter{}, handler)
}

func (w loggingWriter) Write(p []byte) (int, error) {
	n := len(p)
	logs.Logger.Info(string(p[:n-1]))
	return n, nil
}
