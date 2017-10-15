package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

type JSONObject map[string]interface{}

func Respond(w http.ResponseWriter, code int, object JSONObject) {
	w.WriteHeader(code)

	if object != nil {
		bytes, err := json.Marshal(object)
		if err != nil {
			logs.Logger.Panic("Unexpected JSON marshal err", err)
		}

		w.Write(bytes)
	}
}
