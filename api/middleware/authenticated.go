package middleware

import (
	"context"
	"net/http"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

type Validator interface {
	IsValid(authHeader string) (isValid bool, org string)
}

func AuthenticatedHandler(h http.Handler, redirect string, validator Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, org := validator.IsValid(r.Header.Get("Authorization"))
		if !ok {
			http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
			return
		}

		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), persistence.OrgContextKey, org)))
	})
}
