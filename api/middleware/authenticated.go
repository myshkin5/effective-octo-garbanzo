package middleware

import (
	"net/http"
)

type Validator interface {
	IsValid(authHeader string) (isValid bool)
}

func AuthenticatedHandler(h http.Handler, redirect string, validator Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok := validator.IsValid(r.Header.Get("Authorization"))
		if !ok {
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}

		h.ServeHTTP(w, r)
	})
}
