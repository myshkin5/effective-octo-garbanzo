package handlers

import (
	"net/http"
	"sort"
	"strings"
)

// MethodHandler is an http.Handler that dispatches to a handler whose key in the
// MethodHandler's map matches the name of the HTTP request's method, eg: GET
//
// If the request's method is OPTIONS and OPTIONS is not a key in the map then
// the handler responds with a status of 200 and sets the Allow header to a
// comma-separated list of available methods.
//
// If the request's method doesn't match any of its keys the handler responds
// with a status of HTTP 405 "Method Not Allowed" and sets the Allow header to a
// comma-separated list of available methods.
type MethodHandler map[string]http.Handler

func (h MethodHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Copied directly from Gorilla but calls our own Error() func

	if handler, ok := h[req.Method]; ok {
		handler.ServeHTTP(w, req)
	} else {
		allow := []string{}
		for k := range h {
			allow = append(allow, k)
		}
		sort.Strings(allow)
		w.Header().Set("Allow", strings.Join(allow, ", "))
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			Error(w, "Method not allowed", http.StatusMethodNotAllowed, nil, nil)
		}
	}
}
