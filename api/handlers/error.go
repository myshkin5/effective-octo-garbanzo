package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/services"
)

const (
	INVALID_JSON = "Body of request was not valid JSON"
	INVALID_UUID = "Invalid UUID"
)

func Error(w http.ResponseWriter, error string, code int, err error, mapping map[string]string) {
	var validationErrors map[string][]string
	validationError, ok := err.(services.ValidationError)
	if ok {
		code = http.StatusBadRequest
		validationErrors = validationError.Errors()
	}

	message := "Returning %d, message %s"
	messageWithError := message + ", error %v"
	if code >= http.StatusInternalServerError {
		if err != nil {
			logs.Logger.Errorf(messageWithError, code, error, err)
		} else {
			logs.Logger.Errorf(message, code, error)
		}
	} else if code >= http.StatusBadRequest {
		if err != nil {
			logs.Logger.Warnf(messageWithError, code, error, err)
		} else {
			logs.Logger.Warnf(message, code, error)
		}
	} else {
		if err != nil {
			logs.Logger.Infof(messageWithError, code, error, err)
		} else {
			logs.Logger.Infof(message, code, error)
		}
	}

	ret := JSONObject{
		"code":   code,
		"error":  error,
		"status": http.StatusText(code),
	}

	if len(validationErrors) > 0 {
		var errorList []string
		for field, errors := range validationErrors {
			remappedField, ok := mapping[field]
			if !ok {
				remappedField = field
			}
			for _, err := range errors {
				errorList = append(errorList, fmt.Sprintf("%s %s", remappedField, err))
			}
		}
		ret["errors"] = errorList
	}

	bytes, err := json.Marshal(ret)
	if err != nil {
		logs.Logger.Panic("Unexpected JSON marshal err", err)
	}

	w.WriteHeader(code)
	w.Write(bytes)
}
