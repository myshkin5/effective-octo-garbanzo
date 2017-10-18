package middleware

import (
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

func StandardHeadersHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hw := &hijackedWriter{innerWriter: w}
		h.ServeHTTP(hw, r)
	})
}

type hijackedWriter struct {
	innerWriter http.ResponseWriter
	wroteHeader bool
}

func (w *hijackedWriter) Header() http.Header {
	return w.innerWriter.Header()
}

func (w *hijackedWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		logs.Logger.Panic("Call WriteHeader() prior to calling Write(), 200 - Ok is not assumed")
	}

	return w.innerWriter.Write(b)
}

func (w *hijackedWriter) WriteHeader(code int) {
	if code != http.StatusNoContent {
		w.Header().Set("Content-Type", "application/json")
	}

	w.innerWriter.WriteHeader(code)

	w.wroteHeader = true
}
