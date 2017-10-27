package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/services"
)

const (
	INVALID_JSON = "Body of request was not valid JSON"
	INVALID_UUID = "Invalid UUID"
)

func Error(w http.ResponseWriter, error string, code int, err error) {
	var validationErrors []string
	validationError, ok := err.(services.ValidationError)
	if ok {
		code = http.StatusBadRequest
		validationErrors = validationError.Errors()
	}

	if err != nil {
		logs.Logger.Errorf("Returning %d, message %s, error %v", code, error, err)
	} else {
		logs.Logger.Errorf("Returning %d, message %s", code, error)
	}

	ret := JSONObject{
		"code":   code,
		"error":  error,
		"status": http.StatusText(code),
	}

	if len(validationErrors) > 0 {
		ret["errors"] = validationErrors
	}

	bytes, err := json.Marshal(ret)
	if err != nil {
		logs.Logger.Panic("Unexpected JSON marshal err", err)
	}

	w.WriteHeader(code)
	w.Write(bytes)
}
