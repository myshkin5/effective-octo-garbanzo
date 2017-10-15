package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

const (
	INVALID_JSON = "Body of request was not valid JSON"
	INVALID_UUID = "Invalid UUID"
)

func Error(w http.ResponseWriter, error string, code int, errs ...error) {
	switch len(errs) {
	case 0:
		logs.Logger.Errorf("Returning %d, message %s", code, error)
	case 1:
		logs.Logger.Errorf("Returning %d, message %s, error %v", code, error, errs[0])
	default:
		logs.Logger.Panic("Multiple errors not yet supported")
	}

	ret := JSONObject{
		"code":   code,
		"error":  error,
		"status": http.StatusText(code),
	}
	bytes, err := json.Marshal(ret)
	if err != nil {
		logs.Logger.Panic("Unexpected JSON marshal err", err)
	}

	w.WriteHeader(code)
	w.Write(bytes)
}
