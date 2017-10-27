package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

type JSONObject map[string]interface{}

func Respond(w http.ResponseWriter, code int, v interface{}) {
	w.WriteHeader(code)

	if v != nil {
		bytes, err := json.Marshal(v)
		if err != nil {
			logs.Logger.Panic("Unexpected JSON marshal err", err)
		}

		w.Write(bytes)
	}
}
